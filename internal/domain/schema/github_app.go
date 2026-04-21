package schema

import "time"

// GithubAppInstallationCallbackResponse is returned after a successful
// GitHub App installation callback is processed by the backend.
type GithubAppInstallationCallbackResponse struct {
	Message string                        `json:"message"`
	Data    GithubAppInstallationCallback `json:"data"`
}

// GithubAppInstallationCallback contains persisted installation information
// captured from GitHub during callback handling.
type GithubAppInstallationCallback struct {
	InstallationID int64   `json:"installation_id"`
	AccountLogin   *string `json:"account_login,omitempty"`
	AccountID      *int64  `json:"account_id,omitempty"`
	TargetType     *string `json:"target_type,omitempty"`
}

// GithubAppInstallationListResponse is the API response for listing GitHub App
// installations connected to a user.
type GithubAppInstallationListResponse struct {
	Message string                  `json:"message"`
	Data    []GithubAppInstallation `json:"data"`
}

// GithubAppInstallation represents a GitHub App installation record exposed
// by API responses.
type GithubAppInstallation struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	InstallationID int64   `json:"installation_id"`
	AccountLogin   *string `json:"account_login,omitempty"`
	AccountID      *int64  `json:"account_id,omitempty"`
	TargetType     *string `json:"target_type,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      *string `json:"updated_at,omitempty"`
}

// GithubAppCreateInstallationTokenResponse is returned when the backend
// successfully creates a GitHub installation access token.
type GithubAppCreateInstallationTokenResponse struct {
	Message string                            `json:"message"`
	Data    GithubAppInstallationTokenPayload `json:"data"`
}

// GithubAppInstallationTokenPayload contains token details returned to clients
// for short-lived installation access.
type GithubAppInstallationTokenPayload struct {
	Token        string   `json:"token"`
	ExpiresAt    string   `json:"expires_at"`
	Permissions  any      `json:"permissions,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
}

// GithubAppInstallationTokenAPIResponse models GitHub's installation token API
// response payload.
type GithubAppInstallationTokenAPIResponse struct {
	Token                  string   `json:"token"`
	ExpiresAt              string   `json:"expires_at"`
	Permissions            any      `json:"permissions"`
	RepositorySelection    string   `json:"repository_selection"`
	RepositoriesURL        string   `json:"repositories_url"`
	SingleFile             string   `json:"single_file"`
	SingleFilePaths        []string `json:"single_file_paths"`
	HasMultipleSingleFiles bool     `json:"has_multiple_single_files"`
}

// GithubAppInstallURLResponse is returned when the backend generates the
// GitHub App installation URL.
type GithubAppInstallURLResponse struct {
	Message string                  `json:"message"`
	Data    GithubAppInstallURLData `json:"data"`
}

// GithubAppInstallURLData contains the GitHub installation URL that clients
// should redirect users to.
type GithubAppInstallURLData struct {
	URL string `json:"url"`
}

// GithubAppInstallationToken is the response payload for a GitHub
// installation access token request.
type GithubAppInstallationToken struct {
	Token               string    `json:"token"`
	ExpiresAt           time.Time `json:"expires_at"`
	Permissions         any       `json:"permissions,omitempty"`
	RepositorySelection string    `json:"repository_selection,omitempty"`
}

type GithubInstallationResponse struct {
	ID         int64  `json:"id"`
	TargetType string `json:"target_type"`
	Account    struct {
		Login string `json:"login"`
		ID    int64  `json:"id"`
	} `json:"account"`
}

// GithubInstallation represents a persisted GitHub App installation record
// mapped to an internal user.
type GithubInstallation struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	InstallationID int64      `json:"installation_id"`
	AccountLogin   *string    `json:"account_login,omitempty"`
	AccountID      *int64     `json:"account_id,omitempty"`
	TargetType     *string    `json:"target_type,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

// GithubInstallationRepository represents a repository accessible through
// a specific GitHub App installation.
type GithubInstallationRepository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch,omitempty"`
}
