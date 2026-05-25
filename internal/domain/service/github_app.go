package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	github_app_repository "github.com/server-selfish/backend/internal/domain/repository/github_app"
	"github.com/server-selfish/backend/internal/domain/schema"
	cache_infra "github.com/server-selfish/backend/internal/infra/cache"
	github_infra "github.com/server-selfish/backend/internal/infra/github"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
	"github.com/spf13/viper"
	"github.com/valkey-io/valkey-go"
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
		) error
		// ListInstallations returns all GitHub App installations connected to
		// the given user.
		ListInstallations(ctx context.Context, userID pgtype.UUID) ([]schema.GithubInstallation, error)
		// ListInstallationRepositories returns repositories accessible by the
		// provided installation, after validating the installation belongs to user.
		ListInstallationRepositories(ctx context.Context, userID pgtype.UUID, installationID int64) ([]schema.GithubInstallationRepository, error)
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
		installStatePrefix string
		githubInfra        github_infra.GithubInfra
	}
)

// NewGithubAppService constructs a GithubAppService and validates required
// GitHub App configuration from application config.
func NewGithubAppService(repo *github_app_repository.Queries, cache cache_infra.ValkeyInfra, gi github_infra.GithubInfra, logger zerolog.Logger) (GithubAppService, error) {
	appID := strings.TrimSpace(viper.GetString("auth.github.app.id"))
	appSlug := strings.TrimSpace(viper.GetString("auth.github.app.slug"))
	callbackURI := fmt.Sprintf("%s/api/github-app/callback", viper.Get("app.base.url"))
	baseInstallURL := "https://github.com/apps"
	privateKeyPath := viper.GetString("auth.github.private.key.path")

	var privateKeyPEM string

	if privateKeyPath != "" {
		keyData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to read private key file")
		}
		privateKeyPEM = string(keyData)
	}

	if appID == "" || appSlug == "" {
		logger.Fatal().Msg(defined_error.ErrMissingAppIDOrSlug.Error())
	}
	if strings.TrimSpace(privateKeyPEM) == "" {
		logger.Fatal().Msg(defined_error.ErrMissingGithubAppPrivateKey.Error())
	}
	if callbackURI == "" {
		logger.Fatal().Msg(defined_error.ErrMissingGithubAppCallbackURI.Error())
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
		installStatePrefix: "github_app_install_state:",
		githubInfra:        gi,
	}, nil
}

// GetInstallURL generates a GitHub App installation URL and stores a temporary
// state-to-user mapping in cache for callback verification.
func (s *githubAppService) GetInstallURL(ctx context.Context, userID string) (string, string, error) {
	// generate random string just for state passes in query params below
	state, err := pkg.GenerateState(32)
	if err != nil {
		return "", "", err
	}

	// parse to normalize url
	u, err := url.Parse(fmt.Sprintf("%s/%s/installations/new", s.baseInstallURL, s.appSlug))
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", defined_error.ErrParseUrl.Error(), err)
	}
	q := u.Query()
	q.Set("state", state)
	q.Set("redirect_uri", s.callbackURI)
	u.RawQuery = q.Encode()

	// save the state in cache for callback comparison later
	if err := s.cache.Set(ctx, s.installStatePrefix+state, userID, 10*time.Minute); err != nil {
		return "", "", err
	}

	return u.String(), state, nil
}

// HandleInstallCallback validates callback state, resolves the associated user,
// fetches installation details from GitHub, and upserts the mapping.
func (s *githubAppService) HandleInstallCallback(
	ctx context.Context,
	state string,
	installationID int64,
) error {
	// get the user id by state passed
	cacheKey := s.installStatePrefix + state
	userID, err := s.cache.Get(ctx, cacheKey)
	if err != nil || strings.TrimSpace(userID) == "" {
		return err
	}

	// delete the key after retrieved
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		return err
	}

	// installation details
	appInstall, err := s.githubInfra.FetchInstallation(ctx, installationID)
	if err != nil {
		return err
	}

	pgUserID, err := pkg.StringToPgUUID(userID)
	if err != nil {
		return defined_error.ErrStringUUIDTypeCasting
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
		return fmt.Errorf("%s: %w", defined_error.ErrUpsertGithubInstallation.Error(), err)
	}

	return nil
}

// ListInstallations returns all persisted GitHub App installations for the
// specified user.
func (s *githubAppService) ListInstallations(ctx context.Context, userID pgtype.UUID) ([]schema.GithubInstallation, error) {
	// get all installation by User ID
	rows, err := s.repo.ListGithubInstallationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", defined_error.ErrGetListGithubInstallation.Error(), err)
	}

	out := make([]schema.GithubInstallation, 0, len(rows))
	for _, r := range rows {
		item := schema.GithubInstallation{
			InstallationID: r.InstallationID,
		}

		if r.AccountLogin.Valid {
			v := r.AccountLogin.String
			item.AccountLogin = &v
		}

		out = append(out, item)
	}

	return out, nil
}

// ListInstallationRepositories returns repositories accessible by a GitHub App
// installation after ensuring the installation belongs to the user.
func (s *githubAppService) ListInstallationRepositories(ctx context.Context, userID pgtype.UUID, installationID int64) ([]schema.GithubInstallationRepository, error) {

	if _, err := s.repo.GetGithubInstallationByUserIDAndInstallationID(ctx, github_app_repository.GetGithubInstallationByUserIDAndInstallationIDParams{
		UserID:         userID,
		InstallationID: installationID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, defined_error.ErrGetGithubInstallationNotFound
		}
		return nil, defined_error.ErrGetGithubInstallation
	}

	cacheKey := fmt.Sprintf("github:repos:%s:%d", userID, installationID)
	var repositories []schema.GithubInstallationRepository
	getCacheErr := s.cache.GetJSON(ctx, cacheKey, &repositories)

	if getCacheErr != nil {
		if errors.Is(getCacheErr, valkey.Nil) {
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
				return nil, fmt.Errorf("%s: %w", defined_error.ErrListInstallationRepositoriesRequest.Error(), err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Printf("failed to close response body: %v", err)
				}
			}()

			if resp.StatusCode >= 400 {
				return nil, fmt.Errorf("%s: %v", defined_error.ErrListInstallationRepositoriesRequest.Error(), resp)
			}

			var out struct {
				Repositories []schema.GithubInstallationRepository `json:"repositories"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return nil, fmt.Errorf("%s: %w", defined_error.ErrDecodeInstallationRepositoriesResponse.Error(), err)
			}

			for _, repo := range out.Repositories {
				repositories = append(repositories, schema.GithubInstallationRepository{
					ID:       repo.ID,
					Name:     repo.Name,
					FullName: repo.FullName,
				})
			}

			_ = s.cache.SetJSON(ctx, cacheKey, repositories, 5*time.Minute)
		} else {
			return nil, getCacheErr
		}
	}
	return repositories, nil
}
