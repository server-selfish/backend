package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	github_app_repository "github.com/server-selfish/backend/internal/domain/repository/github_app"
	"github.com/server-selfish/backend/internal/domain/schema"
	cache_infra "github.com/server-selfish/backend/internal/infra/cache"
	github_infra "github.com/server-selfish/backend/internal/infra/github"
	"github.com/server-selfish/backend/internal/pkg"
	"github.com/spf13/viper"
)

type (
	// GithubAppService defines business operations for GitHub App installation
	// flows, persisted installation mappings, and repository access via backend.
	GithubAppService interface {
		// GetInstallURL creates a GitHub App installation URL for a user and
		// returns the generated URL together with the CSRF state value.
		GetInstallURL(ctx context.Context, userID string) (string, string, error)
		// HandleInstallCallback validates and consumes the installation callback
		// state, fetches installation metadata from GitHub, and upserts it to DB.
		HandleInstallCallback(
			ctx context.Context,
			state string,
			installationID int64,
			setupAction string,
		) error
		// ListInstallations returns all GitHub App installations connected to
		// the given user.
		ListInstallations(ctx context.Context, userID string) ([]schema.GithubInstallation, error)
		// ListInstallationRepositories returns repositories accessible by the
		// provided installation, after validating the installation belongs to user.
		ListInstallationRepositories(ctx context.Context, userID string, installationID int64) ([]schema.GithubInstallationRepository, error)
	}

	// githubAppService is the concrete implementation of GithubAppService.
	githubAppService struct {
		repo               *github_app_repository.Queries
		cache              cache_infra.ValkeyInfra
		httpClient         *http.Client
		appID              string
		appSlug            string
		privateKeyPEM      string
		callbackURI        string
		baseInstallURL     string
		githubAPIBaseURL   string
		stateTTL           time.Duration
		installStatePrefix string
		githubInfra        github_infra.GithubInfra
	}
)

// NewGithubAppService constructs a GithubAppService and validates required
// GitHub App configuration from application config.
func NewGithubAppService(repo *github_app_repository.Queries, cache cache_infra.ValkeyInfra, gi github_infra.GithubInfra) (GithubAppService, error) {
	appID := strings.TrimSpace(viper.GetString("auth.github.app.id"))
	appSlug := strings.TrimSpace(viper.GetString("auth.github.app.slug"))
	privateKeyPEM := viper.GetString("auth.github.private.key.pem")
	callbackURI := fmt.Sprintf("%s/github-app/callback", viper.Get("app.base.url"))
	baseInstallURL := "https://github.com/apps"

	if appID == "" || appSlug == "" {
		return nil, errors.New("github app configuration is incomplete: app_id and app_slug are required")
	}
	if strings.TrimSpace(privateKeyPEM) == "" {
		return nil, errors.New("github app private key is required (private_key_pem or private_key_path)")
	}
	if callbackURI == "" {
		return nil, errors.New("github app callback_uri is required")
	}

	return &githubAppService{
		repo:               repo,
		cache:              cache,
		httpClient:         &http.Client{Timeout: 20 * time.Second},
		appID:              appID,
		appSlug:            appSlug,
		privateKeyPEM:      privateKeyPEM,
		callbackURI:        callbackURI,
		baseInstallURL:     strings.TrimRight(baseInstallURL, "/"),
		githubAPIBaseURL:   "https://api.github.com",
		stateTTL:           10 * time.Minute,
		installStatePrefix: "github_app_install_state:",
		githubInfra:        gi,
	}, nil
}

// GetInstallURL generates a GitHub App installation URL and stores a temporary
// state-to-user mapping in cache for callback verification.
func (s *githubAppService) GetInstallURL(ctx context.Context, userID string) (string, string, error) {

	if strings.TrimSpace(userID) == "" {
		return "", "", errors.New("user id is required")
	}

	// generate random string just for state passes in query params below
	state, err := pkg.GenerateState(32)
	if err != nil {
		return "", "", fmt.Errorf("generate state: %w", err)
	}

	// parse to normalize url
	u, err := url.Parse(fmt.Sprintf("%s/%s/installations/new", s.baseInstallURL, s.appSlug))
	if err != nil {
		return "", "", err
	}
	q := u.Query()
	q.Set("state", state)
	q.Set("redirect_uri", s.callbackURI)
	u.RawQuery = q.Encode()

	// save the state in cache for callback comparison later
	if err := s.cache.Set(ctx, s.installStatePrefix+state, userID, s.stateTTL); err != nil {
		return "", "", fmt.Errorf("store install state: %w", err)
	}

	return u.String(), state, nil
}

// HandleInstallCallback validates callback state, resolves the associated user,
// fetches installation details from GitHub, and upserts the mapping.
func (s *githubAppService) HandleInstallCallback(
	ctx context.Context,
	state string,
	installationID int64,
	setupAction string,
) error {
	if strings.TrimSpace(state) == "" {
		return errors.New("state is required")
	}
	if installationID <= 0 {
		return errors.New("installation id is required")
	}

	// get the user id by state passed
	cacheKey := s.installStatePrefix + state
	userID, err := s.cache.Get(ctx, cacheKey)
	if err != nil || strings.TrimSpace(userID) == "" {
		return errors.New("invalid or expired state")
	}

	// delete the key after retrieved
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		return fmt.Errorf("consume install state: %w", err)
	}

	// installation details
	appInstall, err := s.githubInfra.FetchInstallation(ctx, installationID)
	if err != nil {
		return err
	}

	pgUserID, err := pkg.StringToPgUUID(userID)
	if err != nil {
		return pkg.ErrBadRequest
	}

	recordID := pkg.PgUUIDFromUUID(uuid.New())

	var accountLogin pgtype.Text
	if strings.TrimSpace(appInstall.Account.Login) != "" {
		accountLogin = pkg.ToPgText(&appInstall.Account.Login)
	}

	var accountID pgtype.Int8
	if appInstall.Account.ID > 0 {
		accountID = pgtype.Int8{Int64: appInstall.Account.ID, Valid: true}
	}

	var targetType pgtype.Text
	if strings.TrimSpace(appInstall.TargetType) != "" {
		targetType = pkg.ToPgText(&appInstall.TargetType)
	}

	// upsert the data to DB
	_, err = s.repo.UpsertGithubInstallation(ctx, github_app_repository.UpsertGithubInstallationParams{
		ID:             recordID,
		UserID:         pgUserID,
		InstallationID: appInstall.ID,
		AccountLogin:   accountLogin,
		AccountID:      accountID,
		TargetType:     targetType,
	})
	if err != nil {
		return fmt.Errorf("upsert github installation: %w", err)
	}

	_ = setupAction
	return nil
}

// ListInstallations returns all persisted GitHub App installations for the
// specified user.
func (s *githubAppService) ListInstallations(ctx context.Context, userID string) ([]schema.GithubInstallation, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	pgUserID, err := pkg.StringToPgUUID(userID)
	if err != nil {
		return nil, pkg.ErrBadRequest
	}

	// get all installation by User ID
	rows, err := s.repo.ListGithubInstallationsByUserID(ctx, pgUserID)
	if err != nil {
		return nil, fmt.Errorf("list installations: %w", err)
	}

	out := make([]schema.GithubInstallation, 0, len(rows))
	for _, r := range rows {
		item := schema.GithubInstallation{
			ID:             r.ID.String(),
			UserID:         r.UserID.String(),
			InstallationID: r.InstallationID,
			CreatedAt:      r.CreatedAt.Time,
		}

		if r.AccountLogin.Valid {
			v := r.AccountLogin.String
			item.AccountLogin = &v
		}
		if r.AccountID.Valid {
			v := r.AccountID.Int64
			item.AccountID = &v
		}
		if r.TargetType.Valid {
			v := r.TargetType.String
			item.TargetType = &v
		}
		if r.UpdatedAt.Valid {
			v := r.UpdatedAt.Time
			item.UpdatedAt = &v
		}

		out = append(out, item)
	}

	return out, nil
}

// ListInstallationRepositories returns repositories accessible by a GitHub App
// installation after ensuring the installation belongs to the user.
func (s *githubAppService) ListInstallationRepositories(ctx context.Context, userID string, installationID int64) ([]schema.GithubInstallationRepository, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}
	if installationID <= 0 {
		return nil, errors.New("installation id is required")
	}

	pgUserID, err := pkg.StringToPgUUID(userID)
	if err != nil {
		return nil, pkg.ErrBadRequest
	}

	_, err = s.repo.GetGithubInstallationByUserIDAndInstallationID(ctx, github_app_repository.GetGithubInstallationByUserIDAndInstallationIDParams{
		UserID:         pgUserID,
		InstallationID: installationID,
	})
	if err != nil {
		return nil, pkg.ErrNotFound
	}

	tokenRes, err := s.githubInfra.CreateInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/installation/repositories?per_page=100", s.githubAPIBaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+tokenRes.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list installation repositories request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("list installation repositories failed: status=%d", resp.StatusCode)
	}

	var out struct {
		Repositories []schema.GithubInstallationRepository `json:"repositories"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode installation repositories response: %w", err)
	}

	repositories := make([]schema.GithubInstallationRepository, 0, len(out.Repositories))
	for _, repo := range out.Repositories {
		repositories = append(repositories, schema.GithubInstallationRepository{
			ID:            repo.ID,
			Name:          repo.Name,
			FullName:      repo.FullName,
			Private:       repo.Private,
			HTMLURL:       repo.HTMLURL,
			CloneURL:      repo.CloneURL,
			DefaultBranch: repo.DefaultBranch,
		})
	}

	return repositories, nil
}
