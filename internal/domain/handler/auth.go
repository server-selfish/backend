package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
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
		as service.AuthService
	}
)

func NewAuthHandler(as service.AuthService) AuthHandler {
	return &authHandler{
		as: as,
	}
}

func (h *authHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {

	// just redirect to login URL with client id and redirect callback URI defined
	ctx := r.Context()

	loginURL, err := h.as.GetGithubLoginURL(ctx)
	if err != nil {
		pkg.WriteJSON(w, http.StatusInternalServerError, schema.AuthErrorResponse{
			Message: "failed to build github login url",
			Error:   err.Error(),
		})
		return
	}

	// Optional helper for frontend if it wants to keep track of state.
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "oauth_state",
	// 	Value:    state,
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
	// 	SameSite: http.SameSiteLaxMode,
	// 	MaxAge:   600,
	// })

	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
}

func (h *authHandler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validation query params
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		pkg.WriteJSON(w, http.StatusBadRequest, schema.AuthErrorResponse{
			Message: "invalid callback params",
			Error:   "missing code or state",
		})
		return
	}

	// handler github callback service call
	tokenPair, err := h.as.HandleGithubCallback(ctx, code, state, r.UserAgent(), pkg.ReadIP(r))
	if err != nil {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "github authentication fai 	led",
			Error:   err.Error(),
		})
		return
	}

	// set refresh token to cookie httponly
	if tokenPair.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokenPair.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.RefreshTokenExpiresIn),
		})
	}

	// send back response with access token
	pkg.WriteJSON(w, http.StatusOK, schema.GithubCallbackResponse{
		Message: "github login success",
		Data: schema.AuthTokenPair{
			AccessToken:           tokenPair.AccessToken,
			AccessTokenExpiresIn:  tokenPair.AccessTokenExpiresIn,
			RefreshTokenExpiresIn: tokenPair.RefreshTokenExpiresIn,
			TokenType:             tokenPair.TokenType,
		},
	})
}

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.RefreshTokenRequest
	// get refresh token from body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// fallback to cookie if body isn't provided
		req.RefreshToken = ""
	}

	// get refresh token from cookie
	if strings.TrimSpace(req.RefreshToken) == "" {
		if c, err := r.Cookie("refresh_token"); err == nil {
			req.RefreshToken = c.Value
		}
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		pkg.WriteJSON(w, http.StatusBadRequest, schema.AuthErrorResponse{
			Message: "refresh token is required",
		})
		return
	}

	// get the new access token inside token pair
	tokenPair, err := h.as.RefreshAccessToken(ctx, req.RefreshToken, r.UserAgent(), pkg.ReadIP(r))
	if err != nil {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "failed to refresh token",
			Error:   err.Error(),
		})
		return
	}

	// rotate cookie refresh token
	if tokenPair.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokenPair.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.RefreshTokenExpiresIn),
		})
	}

	// send back response with access token
	pkg.WriteJSON(w, http.StatusOK, schema.RefreshTokenResponse{
		Message: "token refreshed",
		Data: schema.AuthTokenPair{
			AccessToken:           tokenPair.AccessToken,
			AccessTokenExpiresIn:  tokenPair.AccessTokenExpiresIn,
			RefreshTokenExpiresIn: tokenPair.RefreshTokenExpiresIn,
			TokenType:             tokenPair.TokenType,
		},
	})
}

func (h *authHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user id from context (ingested from middleware)
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	// get user profile
	me, err := h.as.GetMe(ctx, userID)
	if err != nil {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "failed to resolve user",
			Error:   err.Error(),
		})
		return
	}

	pkg.WriteJSON(w, http.StatusOK, schema.MeResponse{
		Message: "ok",
		Data: schema.UserMeData{
			ID:             me.ID,
			Provider:       me.Provider,
			ProviderUserID: me.ProviderUserID,
			Username:       me.Username,
			Email:          me.Email,
			AvatarURL:      me.AvatarURL,
		},
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.RefreshToken = ""
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		if c, err := r.Cookie("refresh_token"); err == nil {
			req.RefreshToken = c.Value
		}
	}
	if strings.TrimSpace(req.RefreshToken) != "" {
		// logout service call
		_ = h.as.Logout(ctx, req.RefreshToken)
	}

	// clear cookie regardless of service result
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	pkg.WriteJSON(w, http.StatusOK, schema.LogoutResponse{
		Message: "logout success",
	})
}
