package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	user_repository "github.com/server-selfish/backend/internal/domain/repository/user"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
	"github.com/spf13/viper"
)

type (
	AuthService interface {
		GetGithubLoginURL(ctx context.Context) (loginURL string, err error)
		HandleGithubCallback(ctx context.Context, code string, userAgent string, ipAddress string) (schema.AuthTokenPair, error)
		RefreshAccessToken(ctx context.Context, refreshToken string, userAgent string, ipAddress string) (schema.AuthTokenPair, error)
		Logout(ctx context.Context, refreshToken string) error
		GetMe(ctx context.Context, userID pgtype.UUID) (schema.UserMeData, error)
		ParseAccessToken(accessToken string) (*pkg.AccessTokenClaims, error)
	}

	authService struct {
		ur           *user_repository.Queries
		token        *pkg.TokenManager
		httpClient   *http.Client
		clientID     string
		clientSecret string
		redirectURI  string
	}
)

func NewAuthService(ur *user_repository.Queries, tm *pkg.TokenManager, logger zerolog.Logger) (AuthService, error) {
	clientID := viper.GetString("auth.github.client.id")
	clientSecret := viper.GetString("auth.github.client.secret")
	redirectURI := fmt.Sprintf("%s/auth/github/callback", viper.GetString("app.base.url"))

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		logger.Fatal().Msg("github oauth configuration is incomplete")
	}

	return &authService{
		ur:           ur,
		token:        tm,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}, nil
}

func (a *authService) GetGithubLoginURL(ctx context.Context) (string, error) {
	u, err := url.Parse("https://github.com/login/oauth/authorize")
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("client_id", a.clientID)
	q.Set("redirect_uri", a.redirectURI)
	q.Set("scope", "read:user user:email")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (a *authService) HandleGithubCallback(ctx context.Context, code string, userAgent string, ipAddress string) (schema.AuthTokenPair, error) {
	// get the github access token to be accessed by request
	ghToken, err := a.exchangeGithubCode(ctx, code)
	if err != nil {
		return schema.AuthTokenPair{}, err
	}

	// get the github profile using prev retrieved access token
	profile, err := a.fetchGithubProfile(ctx, ghToken)
	if err != nil {
		return schema.AuthTokenPair{}, err
	}

	if strings.TrimSpace(profile.Login) == "" || profile.ID == 0 {
		return schema.AuthTokenPair{}, errors.New("invalid github profile response")
	}

	userUUID := pkg.PgUUIDFromUUID(uuid.New())

	// upsert user data to DB
	user, err := a.ur.UpsertGithubUser(ctx, user_repository.UpsertGithubUserParams{
		ID:             userUUID,
		ProviderUserID: profile.ID,
		Username:       profile.Login,
		Email:          pkg.ToPgText(profile.Email),
		AvatarUrl:      pkg.ToPgText(profile.AvatarURL),
	})
	if err != nil {
		return schema.AuthTokenPair{}, fmt.Errorf("%s: %w", defined_error.ErrUpsertGithubUser.Error(), err)
	}

	return a.issueSessionTokenPair(ctx, user.ID, userAgent, ipAddress)
}

func (a *authService) RefreshAccessToken(ctx context.Context, refreshToken string, userAgent string, ipAddress string) (schema.AuthTokenPair, error) {
	// hash the refresh token from client
	refreshHash := pkg.HashToken(refreshToken)

	// get user session based on hashed refresh token (this also filter expires_at > now())
	session, err := a.ur.GetAuthSessionByRefreshTokenHash(ctx, refreshHash)
	if err != nil {
		return schema.AuthTokenPair{}, defined_error.ErrInvalidRefreshTokenInSession
	}

	// generate new refresh token for new session (session = as long as particular access_token is valid)
	// newRawRefresh, newRefreshHash, newRefreshExp, err := a.token.NewRefreshToken()
	// if err != nil {
	// 	return schema.AuthTokenPair{}, err
	// }

	// save new refresh token in session DB
	// if err := a.ur.RotateAuthSessionToken(ctx, user_repository.RotateAuthSessionTokenParams{
	// 	ID:               session.ID,
	// 	RefreshTokenHash: newRefreshHash,
	// 	ExpiresAt:        pkg.PgTimestamptz(newRefreshExp),
	// }); err != nil {
	// 	return schema.AuthTokenPair{}, fmt.Errorf("%s: %w", defined_error.ErrRotateRefreshToken.Error(), err)
	// }

	// generate access token
	access, accessExp, err := a.token.GenerateAccessToken(session.UserID.String(), session.ID.String(), "github")
	if err != nil {
		return schema.AuthTokenPair{}, fmt.Errorf("%s: %w", defined_error.ErrGenerateAccessToken.Error(), err)
	}

	return schema.AuthTokenPair{
		AccessToken:          access,
		AccessTokenExpiresIn: int64(time.Until(accessExp).Seconds()),
		TokenType:            "Bearer",
	}, nil
}

func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	// get auth session
	session, err := a.ur.GetAuthSessionByRefreshTokenHash(ctx, pkg.HashToken(refreshToken))
	if err != nil {
		return fmt.Errorf("%s: %w", defined_error.ErrGetAuthSessionByRefreshToken.Error(), err)
	}

	// update the session to revoked
	if err := a.ur.RevokeAuthSessionByID(ctx, session.ID); err != nil {
		return fmt.Errorf("%s: %w", defined_error.ErrRevokeSession.Error(), err)
	}

	return nil
}

func (a *authService) GetMe(ctx context.Context, userID pgtype.UUID) (schema.UserMeData, error) {
	u, err := a.ur.GetUserByID(ctx, userID)
	if err != nil {
		return schema.UserMeData{}, fmt.Errorf("%s: %w", defined_error.ErrGetUserData.Error(), err)
	}

	var email *string
	if u.Email.Valid {
		v := u.Email.String
		email = &v
	}

	var avatar *string
	if u.AvatarUrl.Valid {
		v := u.AvatarUrl.String
		avatar = &v
	}

	return schema.UserMeData{
		ID:             u.ID.String(),
		Provider:       u.Provider,
		ProviderUserID: u.ProviderUserID,
		Username:       u.Username,
		Email:          email,
		AvatarURL:      avatar,
	}, nil
}

func (a *authService) ParseAccessToken(accessToken string) (*pkg.AccessTokenClaims, error) {
	// parse access token
	return a.token.ParseAccessToken(accessToken)
}

func (a *authService) issueSessionTokenPair(ctx context.Context, userID pgtype.UUID, userAgent string, ipAddress string) (schema.AuthTokenPair, error) {
	sessionID := pkg.PgUUIDFromUUID(uuid.New())

	// create new refresh token
	rawRefresh, refreshHash, refreshExp, err := a.token.NewRefreshToken()
	if err != nil {
		return schema.AuthTokenPair{}, err
	}

	// insert hashed refresh token to DB
	_, err = a.ur.CreateAuthSession(ctx, user_repository.CreateAuthSessionParams{
		ID:               sessionID,
		UserID:           userID,
		RefreshTokenHash: refreshHash,
		ExpiresAt:        pkg.PgTimestamptz(refreshExp),
		UserAgent:        pkg.ToPgText(pkg.StrPtr(userAgent)),
		IpAddress:        pkg.ToPgText(pkg.StrPtr(ipAddress)),
	})
	if err != nil {
		return schema.AuthTokenPair{}, fmt.Errorf("%s: %w", defined_error.ErrCreateAuthSession.Error(), err)
	}

	// generate access token
	access, accessExp, err := a.token.GenerateAccessToken(userID.String(), sessionID.String(), "github")
	if err != nil {
		return schema.AuthTokenPair{}, fmt.Errorf("%s: %w", defined_error.ErrGenerateAccessToken.Error(), err)
	}

	// return all token data
	return schema.AuthTokenPair{
		AccessToken:           access,
		AccessTokenExpiresIn:  int64(time.Until(accessExp).Seconds()),
		RefreshToken:          rawRefresh,
		RefreshTokenExpiresIn: int64(time.Until(refreshExp).Seconds()),
		TokenType:             "Bearer",
	}, nil
}

func (a *authService) exchangeGithubCode(ctx context.Context, code string) (string, error) {
	values := url.Values{}
	values.Set("client_id", a.clientID)
	values.Set("client_secret", a.clientSecret)
	values.Set("code", code)
	values.Set("redirect_uri", a.redirectURI)

	// request builder to get github access token with client_id, client_secret, code, redirectURI in body
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", defined_error.ErrGithubAccessTokenRequestFailed.Error(), err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("%s: %v", defined_error.ErrGithubTokenExchangeFailed.Error(), resp)
	}

	var tokenRes schema.GithubOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return "", fmt.Errorf("%s: %w", defined_error.ErrDecodeGithubTokenResponse.Error(), err)
	}
	if tokenRes.AccessToken == "" {
		return "", defined_error.ErrEmptyGithubAccessToken
	}

	return tokenRes.AccessToken, nil
}

func (a *authService) fetchGithubProfile(ctx context.Context, githubAccessToken string) (schema.GithubUserProfile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return schema.GithubUserProfile{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+githubAccessToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return schema.GithubUserProfile{}, fmt.Errorf("%s: %w", defined_error.ErrGithubProfileRequestFailed.Error(), err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return schema.GithubUserProfile{}, fmt.Errorf("%s: %v", defined_error.ErrGithubProfileRequestFailed.Error(), resp)
	}

	var profile schema.GithubUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return schema.GithubUserProfile{}, fmt.Errorf("%s: %w", defined_error.ErrDecodeGithubProfileResponse.Error(), err)
	}

	return profile, nil
}
