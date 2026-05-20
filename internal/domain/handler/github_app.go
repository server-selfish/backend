package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
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
		gs service.GithubAppService
	}
)

// NewGithubAppHandler constructs a GitHub App HTTP handler with the required
// GitHub App service dependency.
func NewGithubAppHandler(gs service.GithubAppService) GithubAppHandler {
	return &githubAppHandler{gs: gs}
}

// Install validates the authenticated user, builds a GitHub App install URL,
// and redirects the client to GitHub to continue the installation flow.
func (h *githubAppHandler) Install(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user id to be used later for save installation id
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok || strings.TrimSpace(userID) == "" {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	// get the url and random state to be matched later
	installURL, _, err := h.gs.GetInstallURL(ctx, userID)
	if err != nil {
		pkg.WriteJSON(w, http.StatusInternalServerError, schema.AuthErrorResponse{
			Message: "failed to build github app install url",
			Error:   err.Error(),
		})
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "link generated", struct {
		Link string `json:"link"`
	}{
		Link: installURL,
	})

	// TODO: change the redirect to give back url
	// http.Redirect(w, r, installURL, http.StatusTemporaryRedirect)
}

// Callback handles GitHub's post-installation redirect by validating callback
// parameters and delegating installation persistence to the service layer.
func (h *githubAppHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	installationIDRaw := strings.TrimSpace(r.URL.Query().Get("installation_id"))
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	setupAction := strings.TrimSpace(r.URL.Query().Get("setup_action"))

	if installationIDRaw == "" || state == "" {
		pkg.WriteJSON(w, http.StatusBadRequest, schema.AuthErrorResponse{
			Message: "invalid callback params",
			Error:   "missing installation_id or state",
		})
		return
	}

	// parse installationID to integer
	installationID, err := strconv.ParseInt(installationIDRaw, 10, 64)
	if err != nil || installationID <= 0 {
		pkg.WriteJSON(w, http.StatusBadRequest, schema.AuthErrorResponse{
			Message: "invalid installation id",
			Error:   "installation_id must be positive integer",
		})
		return
	}

	if err := h.gs.HandleInstallCallback(
		ctx,
		state,
		installationID,
		setupAction,
	); err != nil {
		pkg.WriteJSON(w, http.StatusInternalServerError, schema.AuthErrorResponse{
			Message: "failed to save github app installation",
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := githubCallbackTmpl.Execute(w, struct{}{}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ListInstallations returns all GitHub App installations associated with the
// currently authenticated user.
func (h *githubAppHandler) ListInstallations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok || strings.TrimSpace(userID) == "" {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	items, err := h.gs.ListInstallations(ctx, userID)
	if err != nil {
		pkg.WriteJSON(w, http.StatusInternalServerError, schema.AuthErrorResponse{
			Message: "failed to list github app installations",
			Error:   err.Error(),
		})
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "ok",
		"data":    items,
	})
}

// ListInstallationRepositories returns repositories accessible by a specific
// GitHub App installation owned by the authenticated user.
func (h *githubAppHandler) ListInstallationRepositories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok || strings.TrimSpace(userID) == "" {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	installationIDRaw := strings.TrimSpace(chi.URLParam(r, "id"))
	installationID, err := strconv.ParseInt(installationIDRaw, 10, 64)
	if err != nil || installationID <= 0 {
		pkg.WriteJSON(w, http.StatusBadRequest, schema.AuthErrorResponse{
			Message: "invalid installation id",
			Error:   "path param must be positive integer",
		})
		return
	}

	repos, err := h.gs.ListInstallationRepositories(ctx, userID, installationID)
	if err != nil {
		pkg.WriteJSON(w, http.StatusInternalServerError, schema.AuthErrorResponse{
			Message: "failed to list installation repositories",
			Error:   err.Error(),
		})
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "ok",
		"data":    repos,
	})
}
