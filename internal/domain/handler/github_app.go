package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	// GithubAppHandler defines HTTP handlers for GitHub App installation lifecycle
	// operations, including install initiation, callback handling, installation
	// listing, and installation repository listing.
	GithubAppHandler interface {
		// Install starts the GitHub App installation flow for the authenticated user
		// by generating a stateful install URL and redirecting to GitHub.
		Install(w http.ResponseWriter, r *http.Request)
		// Callback completes the installation flow by validating callback params and
		// persisting the GitHub installation mapping to the current user context.
		Callback(w http.ResponseWriter, r *http.Request)
		// ListInstallations returns all GitHub App installations associated with the
		// authenticated user.
		ListInstallations(w http.ResponseWriter, r *http.Request)
		// ListInstallationRepositories returns repositories accessible by a specific
		// GitHub App installation owned by the authenticated user.
		ListInstallationRepositories(w http.ResponseWriter, r *http.Request)
	}

	// githubAppHandler is the concrete implementation of GithubAppHandler.
	githubAppHandler struct {
		gs     service.GithubAppService
		logger *zerolog.Logger
	}
)

// NewGithubAppHandler constructs a GitHub App HTTP handler with the required
// GitHub App service dependency.
func NewGithubAppHandler(gs service.GithubAppService, logger zerolog.Logger) GithubAppHandler {
	return &githubAppHandler{gs: gs, logger: &logger}
}

// Install validates the authenticated user, builds a GitHub App install URL,
// and redirects the client to GitHub to continue the installation flow.
func (h *githubAppHandler) Install(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user id to be used later for save installation id
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok || strings.TrimSpace(userID) == "" {
		h.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	// get the url and random state to be matched later
	installURL, _, err := h.gs.GetInstallURL(ctx, userID)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "link generated", struct {
		Link string `json:"link"`
	}{
		Link: installURL,
	})
}

// Callback handles GitHub's post-installation redirect by validating callback
// parameters and delegating installation persistence to the service layer.
func (h *githubAppHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	installationIDRaw := strings.TrimSpace(r.URL.Query().Get("installation_id"))
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	// setupAction := strings.TrimSpace(r.URL.Query().Get("setup_action"))

	if installationIDRaw == "" || state == "" {
		h.logger.Error().Msg(defined_error.ErrMissingInstallationIDOrStateParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidCallbackParams)
		return
	}

	// parse installationID to integer
	installationID, err := strconv.ParseInt(installationIDRaw, 10, 64)
	if err != nil || installationID <= 0 {
		h.logger.Error().Err(err).Msg(defined_error.ErrInvalidInstallationId.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidInstallationId)
		return
	}

	if err := h.gs.HandleInstallCallback(
		ctx,
		state,
		installationID,
	); err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := githubCallbackTmpl.Execute(w, struct{}{}); err != nil {
		h.logger.Error().Err(err).Msg(defined_error.ErrExecuteTemplateError.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
	}
}

// ListInstallations returns all GitHub App installations associated with the
// currently authenticated user.
func (h *githubAppHandler) ListInstallations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok || strings.TrimSpace(userID) == "" {
		h.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	pgUserID, err := pkg.StringToPgUUID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	items, err := h.gs.ListInstallations(ctx, pgUserID)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch installation list success", items)
}

// ListInstallationRepositories returns repositories accessible by a specific
// GitHub App installation owned by the authenticated user.
func (h *githubAppHandler) ListInstallationRepositories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok || strings.TrimSpace(userID) == "" {
		h.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	installationIDRaw := strings.TrimSpace(chi.URLParam(r, "id"))
	installationID, err := strconv.ParseInt(installationIDRaw, 10, 64)
	if err != nil || installationID <= 0 {
		h.logger.Error().Err(err).Msg(defined_error.ErrInvalidInstallationId.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidInstallationId)
		return
	}

	pgUserID, err := pkg.StringToPgUUID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	repos, err := h.gs.ListInstallationRepositories(ctx, pgUserID, installationID)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch repository list success", repos)
}
