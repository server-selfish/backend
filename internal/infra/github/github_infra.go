package github_infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
	"github.com/spf13/viper"
)

type ()

type (
	GithubInfra interface {
		CreateInstallationToken(ctx context.Context, installationID int64) (schema.GithubAppInstallationToken, error)
		FetchInstallation(ctx context.Context, installationID int64) (schema.GithubInstallationResponse, error)
	}
	githubInfra struct {
		logger           *zerolog.Logger
		githubAPIBaseURL string
		appID            string
		privateKeyPEM    string
	}
)

func NewGithubInfra(logger zerolog.Logger) GithubInfra {
	appID := strings.TrimSpace(viper.GetString("auth.github.app.id"))
	privateKeyPath := viper.GetString("auth.github.private.key.path")

	var privateKeyPEM string

	if privateKeyPath != "" {
		keyData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to read private key file")
		}
		privateKeyPEM = string(keyData)
	}
	if strings.TrimSpace(privateKeyPEM) == "" {
		logger.Fatal().Msg(defined_error.ErrMissingGithubAppPrivateKey.Error())
	}
	return &githubInfra{
		logger:           &logger,
		privateKeyPEM:    privateKeyPEM,
		appID:            appID,
		githubAPIBaseURL: "https://api.github.com",
	}
}

func (g *githubInfra) CreateInstallationToken(ctx context.Context, installationID int64) (schema.GithubAppInstallationToken, error) {
	jwtToken, err := pkg.GenerateAppJWTFromPEM(g.privateKeyPEM, g.appID)
	if err != nil {
		g.logger.Error().Err(err).Msg("error generate app jwt")
		return schema.GithubAppInstallationToken{}, err
	}

	reqURL := fmt.Sprintf("%s/app/installations/%d/access_tokens", g.githubAPIBaseURL, installationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(`{}`))
	if err != nil {
		g.logger.Error().Err(err).Msg("error create request with context instance")
		return schema.GithubAppInstallationToken{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28") // 2026-03-10
	req.Header.Set("Content-Type", "application/json")
	httpClient := http.Client{Timeout: 20 * time.Second}

	resp, err := httpClient.Do(req)
	if err != nil {
		g.logger.Error().Err(err).Msg("create installation token request failed")
		return schema.GithubAppInstallationToken{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode >= 400 {
		g.logger.Error().Str("resp", strconv.Itoa(resp.StatusCode)).Msg("http request got client error")
		return schema.GithubAppInstallationToken{}, fmt.Errorf("create installation token failed: status=%d", resp.StatusCode)
	}

	var out struct {
		Token               string    `json:"token"`
		ExpiresAt           time.Time `json:"expires_at"`
		Permissions         any       `json:"permissions"`
		RepositorySelection string    `json:"repository_selection"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		g.logger.Error().Err(err).Msg("decode installation token response error")
		return schema.GithubAppInstallationToken{}, err
	}
	if strings.TrimSpace(out.Token) == "" {
		g.logger.Error().Msg("empty installation token")
		return schema.GithubAppInstallationToken{}, errors.New("empty installation token")
	}

	return schema.GithubAppInstallationToken{
		Token:               out.Token,
		ExpiresAt:           out.ExpiresAt,
		Permissions:         out.Permissions,
		RepositorySelection: out.RepositorySelection,
	}, nil
}

func (g *githubInfra) FetchInstallation(ctx context.Context, installationID int64) (schema.GithubInstallationResponse, error) {
	jwtToken, err := pkg.GenerateAppJWTFromPEM(g.privateKeyPEM, g.appID)
	if err != nil {
		g.logger.Error().Err(err).Msg("generate app jwt from pem error")
		return schema.GithubInstallationResponse{}, err
	}

	reqURL := fmt.Sprintf("%s/app/installations/%d", g.githubAPIBaseURL, installationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		g.logger.Error().Err(err).Msg("failed create request with context")
		return schema.GithubInstallationResponse{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	httpClient := http.Client{Timeout: 20 * time.Second}

	resp, err := httpClient.Do(req)
	if err != nil {
		g.logger.Error().Err(err).Msg("fetch installation request failed")
		return schema.GithubInstallationResponse{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode >= 400 {
		g.logger.Error().Str("code", strconv.Itoa(resp.StatusCode)).Msg("http request got client error")
		return schema.GithubInstallationResponse{}, fmt.Errorf("fetch installation failed: status=%d", resp.StatusCode)
	}

	var out schema.GithubInstallationResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		g.logger.Error().Err(err).Msg("decode installation response error")
		return schema.GithubInstallationResponse{}, err
	}

	return out, nil
}
