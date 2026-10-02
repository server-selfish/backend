package handler

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
	"github.com/server-selfish/backend/internal/presentation"
	"golang.org/x/sync/singleflight"
)

type (
	AuthHandler interface {
		GithubLogin(w http.ResponseWriter, r *http.Request)
		GithubCallback(w http.ResponseWriter, r *http.Request)
		Refresh(w http.ResponseWriter, r *http.Request)
		Me(w http.ResponseWriter, r *http.Request)
		Logout(w http.ResponseWriter, r *http.Request)
	}

	authHandler struct {
		as     service.AuthService
		logger *zerolog.Logger
	}
)

var (
	githubCallbackTmpl = template.Must(
		template.New("github_callback.html").ParseFS(
			presentation.PresentationEmbed,
			"templates/github_callback.html",
		),
	)
	refreshGroup singleflight.Group
)

func NewAuthHandler(as service.AuthService, logger zerolog.Logger) AuthHandler {
	return &authHandler{
		as:     as,
		logger: &logger,
	}
}

// GithubLogin implements [AuthHandler].
// @Summary      Redirect to GitHub OAuth login
// @Description  Redirects the client to the GitHub OAuth authorization URL.
// @Tags         auth
// @Success      307 "Temporary Redirect"
// @Failure      500 {object} pkg.Response{error=string}
// @Header       307 {string} Location "GitHub OAuth authorization URL"
// @Router       /auth/github/login [get]
func (h *authHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	loginURL, err := h.as.GetGithubLoginURL(ctx)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
}

// GithubCallback implements [AuthHandler].
// GithubCallback godoc
// @Summary      Handle GitHub OAuth callback
// @Description  Handles the GitHub OAuth callback, sets authentication cookies, and returns an HTML response.
// @Tags         auth
// @Produce      html
// @Param        code query string true "GitHub OAuth authorization code"
// @Success      200 {string} string "HTML callback page"
// @Failure      400 {object} pkg.Response{error=string}
// @Failure      401 {object} pkg.Response{error=string}
// @Failure      500 {object} pkg.Response{error=string}
// @Router       /auth/github/callback [get]
func (h *authHandler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Error().Msg(defined_error.ErrMissingCodeInQueryParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidCallbackParams)
		return
	}

	tokenPair, err := h.as.HandleGithubCallback(ctx, code, r.UserAgent(), pkg.ReadIP(r))
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrGithubAuthenticationFailed)
		return
	}

	if tokenPair.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "selfish_refresh_token",
			Value:    tokenPair.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.RefreshTokenExpiresIn),
		})
	}
	if tokenPair.AccessToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "selfish_access_token",
			Value:    tokenPair.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.AccessTokenExpiresIn),
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := githubCallbackTmpl.Execute(w, struct{}{}); err != nil {
		h.logger.Error().Err(err).Msg(defined_error.ErrExecuteTemplateError.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
	}
}

// Refresh implements [AuthHandler].
// Refresh 			godoc
// @Summary     Handles refresh access token
// @Description Refreshes the access token using a refresh token supplied either in the request body or the selfish_refresh_token cookie.
// @Tags        auth
// @Accept 			json
// @Produce     json
// @Param       request body schema.RefreshTokenRequest false "Refresh token request"
// @Success     200 {object} pkg.Response{message="token refreshed"}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /auth/refresh [post]
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, _, _, ok := pkg.DecodeAndValidateBody[schema.RefreshTokenRequest](w, r, h.logger)
	if !ok {
		req.RefreshToken = ""
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		if c, err := r.Cookie("selfish_refresh_token"); err == nil {
			req.RefreshToken = c.Value
		}
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		h.logger.Error().Msg(defined_error.ErrRefreshTokenRequired.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrRefreshTokenRequired)
		return
	}
	// get the new access token inside token pair
	v, err, _ := refreshGroup.Do(req.RefreshToken, func() (any, error) {
		return h.as.RefreshAccessToken(ctx, req.RefreshToken, r.UserAgent(), pkg.ReadIP(r))
	})

	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrFailedToRefreshToken)
		return
	}

	tokenPair, ok := v.(schema.AuthTokenPair)
	if !ok {
		h.logger.Error().Msg(defined_error.ErrFailedCastTokenPairs.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrFailedCastTokenPairs)
		return
	}

	if tokenPair.AccessToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "selfish_access_token",
			Value:    tokenPair.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.AccessTokenExpiresIn),
		})
	}
	pkg.ReturnSuccess(w, http.StatusOK, "token refreshed", nil)
}

// Me implements [AuthHandler].
// Me 					godoc
// @Summary     Handles get user profile
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.UserMeData}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /auth/me [get]
func (h *authHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		h.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	// get user profile
	me, err := h.as.GetMe(ctx, ui)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrFailedToGetUser)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch user success", me)
}

// Logout implements [AuthHandler].
// Logout 			godoc
// @Summary     Handles logout session
// @Description Logout using a refresh token supplied either in the request body or the selfish_refresh_token cookie.
// @Tags        auth
// @Accept 			json
// @Produce     json
// @Param       request body schema.RefreshTokenRequest false "Refresh token"
// @Success     200 {object} pkg.Response{message="logout success"}
// @Router      /auth/logout [post]
func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, _, _, ok := pkg.DecodeAndValidateBody[schema.RefreshTokenRequest](w, r, h.logger)
	if !ok {
		req.RefreshToken = ""
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		if c, err := r.Cookie("selfish_refresh_token"); err == nil {
			req.RefreshToken = c.Value
		}
	}
	if strings.TrimSpace(req.RefreshToken) != "" {
		if err := h.as.Logout(ctx, req.RefreshToken); err != nil {
			h.logger.Error().Msg(err.Error())
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "selfish_refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "selfish_access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	pkg.ReturnSuccess(w, http.StatusOK, "logout success", nil)
}
