package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/service/deployment"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
	"github.com/server-selfish/backend/internal/pkg/webhook"
	"github.com/spf13/viper"
)

const (
	// maxWebhookBodyBytes matches GitHub's documented 25 MB payload cap, above
	// which GitHub does not deliver the event at all.
	maxWebhookBodyBytes = 25 << 20
	// branchRefPrefix distinguishes branch pushes from tag pushes.
	branchRefPrefix = "refs/heads/"
	// defaultWebhookBranch applies when github.webhook.branch is unset.
	defaultWebhookBranch = "main"
)

type (
	WebhookHandler interface {
		RedeployOnPush(w http.ResponseWriter, r *http.Request)
	}
	webhookHandler struct {
		ds     service.DeploymentService
		logger *zerolog.Logger
	}
)

// githubPushPayload is a deliberately narrow slice of GitHub's push event. The
// full event carries commits, compare URLs, and head repo details this handler
// never reads, and decoding it whole would cost memory on every delivery.
type githubPushPayload struct {
	Ref     string `json:"ref"`
	Deleted bool   `json:"deleted"`
	Sender  struct {
		Login string `json:"login"`
	} `json:"sender"`
	Repository struct {
		ID int64 `json:"id"`
	} `json:"repository"`
	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`
}

func NewWebhookHandler(ds service.DeploymentService, logger zerolog.Logger) WebhookHandler {
	return &webhookHandler{
		ds:     ds,
		logger: &logger,
	}
}

// RedeployOnPush implements [WebhookHandler].
//
// Public by design: a GitHub delivery carries no JWT, so it is authenticated
// solely by the HMAC signature over the raw body. The request supplies no user
// identity, only the repository and installation it came from; the deployment
// owner is resolved server-side from that pair, so a caller can never choose
// whose deployment gets rebuilt.
func (h *webhookHandler) RedeployOnPush(w http.ResponseWriter, r *http.Request) {
	secret := viper.GetString("github.webhook.secret")
	if secret == "" {
		h.logger.Error().Msg(defined_error.ErrWebhookDisabled.Error())
		pkg.ReturnError(w, http.StatusNotFound, defined_error.ErrWebhookDisabled)
		return
	}

	deliveryID := r.Header.Get("X-GitHub-Delivery")

	// The raw body is needed for the signature, so read it once and reuse
	// those exact bytes for verification and decoding. Any middleware that
	// decodes the body first, or any re-encoding, invalidates the comparison.
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxWebhookBodyBytes))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to read webhook body")
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidWebhookPayload)
		return
	}

	if err := webhook.VerifyGitHubSignature(secret, body, r.Header.Get("X-Hub-Signature-256")); err != nil {
		h.logger.Warn().Str("delivery_id", deliveryID).Msg(err.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, err)
		return
	}

	// A GitHub App webhook receives every subscribed event at one URL. Only
	// pushes redeploy; the rest are acknowledged so they are not retried.
	if event := r.Header.Get("X-GitHub-Event"); event != "push" {
		pkg.ReturnSuccess(w, http.StatusAccepted, "ignored event", map[string]string{"event": event})
		return
	}

	var payload githubPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.Error().Err(err).Str("delivery_id", deliveryID).Msg("failed to decode webhook payload")
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidWebhookPayload)
		return
	}

	// Tag pushes and branch deletions also arrive as push events, and a deleted
	// branch has no HEAD to clone, so both are acknowledged and dropped.
	branch := viper.GetString("github.webhook.branch")
	if branch == "" {
		branch = defaultWebhookBranch
	}
	if payload.Deleted || payload.Ref != branchRefPrefix+branch {
		pkg.ReturnSuccess(w, http.StatusAccepted, "ignored ref", map[string]string{"ref": payload.Ref})
		return
	}

	h.logger.Info().
		Str("delivery_id", deliveryID).
		Str("sender", payload.Sender.Login).
		Int64("repository_id", payload.Repository.ID).
		Int64("installation_id", payload.Installation.ID).
		Str("ref", payload.Ref).
		Msg("webhook push received")

	// Only the repository lookup and lock claim happen here, both indexed and
	// sub-millisecond, so responding before the build finishes still leaves
	// well inside GitHub's delivery deadline. The builds themselves are already
	// running in background goroutines owned by the service.
	results, err := h.ds.TriggerWebhookRedeploy(r.Context(), int32(payload.Repository.ID), payload.Installation.ID)
	if err != nil {
		if errors.Is(err, defined_error.ErrWebhookNotFound) {
			// Expected whenever a push arrives for an installed repository that
			// is not deployed here. 404 rather than 403 so an unrelated
			// repository cannot probe which deployments exist.
			pkg.ReturnError(w, http.StatusNotFound, err)
			return
		}
		h.logger.Error().Err(err).
			Str("delivery_id", deliveryID).
			Int64("repository_id", payload.Repository.ID).
			Msg("failed to trigger webhook redeploy")
		pkg.ReturnError(w, http.StatusInternalServerError, err)
		return
	}

	pkg.ReturnSuccess(w, http.StatusAccepted, "webhook accepted", map[string]any{
		"repository_id": payload.Repository.ID,
		"ref":           payload.Ref,
		"deployments":   results,
	})
}
