package github_infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/pkg"
	"github.com/spf13/viper"
)

type ()

type (
	GithubInfra interface {
		CreateInstallationToken(ctx context.Context, installationID int64) (schema.GithubAppInstallationToken, error)
		FetchInstallation(ctx context.Context, installationID int64) (schema.GithubInstallationResponse, error)
	}
	githubInfra struct {
		githubAPIBaseURL string
		appID            string
		privateKeyPEM    string
	}
)

func NewGithubInfra() GithubInfra {
	appID := strings.TrimSpace(viper.GetString("auth.github.app.id"))
	privateKeyPEM := viper.GetString("auth.github.private.key.pem")

	return &githubInfra{
		privateKeyPEM:    privateKeyPEM,
		appID:            appID,
		githubAPIBaseURL: "https://api.github.com",
	}
}

func (g *githubInfra) CreateInstallationToken(ctx context.Context, installationID int64) (schema.GithubAppInstallationToken, error) {
	if installationID <= 0 {
		return schema.GithubAppInstallationToken{}, errors.New("installation id is required")
	}

	jwtToken, err := pkg.GenerateAppJWTFromPEM(g.privateKeyPEM, g.appID)
	if err != nil {
		return schema.GithubAppInstallationToken{}, err
	}

	reqURL := fmt.Sprintf("%s/app/installations/%d/access_tokens", g.githubAPIBaseURL, installationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(`{}`))
	if err != nil {
		return schema.GithubAppInstallationToken{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{Timeout: 20 * time.Second}

	resp, err := httpClient.Do(req)
	if err != nil {
		return schema.GithubAppInstallationToken{}, fmt.Errorf("create installation token request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return schema.GithubAppInstallationToken{}, fmt.Errorf("create installation token failed: status=%d", resp.StatusCode)
	}

	var out struct {
		Token               string    `json:"token"`
		ExpiresAt           time.Time `json:"expires_at"`
		Permissions         any       `json:"permissions"`
		RepositorySelection string    `json:"repository_selection"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return schema.GithubAppInstallationToken{}, fmt.Errorf("decode installation token response: %w", err)
	}
	if strings.TrimSpace(out.Token) == "" {
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
		return schema.GithubInstallationResponse{}, err
	}

	reqURL := fmt.Sprintf("%s/app/installations/%d", g.githubAPIBaseURL, installationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return schema.GithubInstallationResponse{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	httpClient := http.Client{Timeout: 20 * time.Second}

	resp, err := httpClient.Do(req)
	if err != nil {
		return schema.GithubInstallationResponse{}, fmt.Errorf("fetch installation request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return schema.GithubInstallationResponse{}, fmt.Errorf("fetch installation failed: status=%d", resp.StatusCode)
	}

	var out schema.GithubInstallationResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return schema.GithubInstallationResponse{}, fmt.Errorf("decode installation response: %w", err)
	}

	return out, nil
}
