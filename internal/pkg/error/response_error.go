package defined_error

import "errors"

// common
var (
	ErrBadRequest            = errors.New("invalid request body")
	ErrInternalServerError   = errors.New("internal server error")
	ErrAlreadyExist          = errors.New("resource is already exist")
	ErrNotFound              = errors.New("resource not found")
	ErrInvalidInstallationId = errors.New("invalid installation ID")
)

// auth
var (
	ErrUnauthorized               = errors.New("unauthorized")
	ErrRefreshTokenRequired       = errors.New("refresh token is required")
	ErrFailedToRefreshToken       = errors.New("failed to refresh token")
	ErrFailedCastTokenPairs       = errors.New("failed to cast token pair")
	ErrFailedToGetUser            = errors.New("failed to get user data")
	ErrGithubAuthenticationFailed = errors.New("github authentication failed")
	ErrInvalidCallbackParams      = errors.New("invalid callback params")
	ErrMissingAccessToken         = errors.New("missing access token")
	ErrInvaliAuthorizationFormat  = errors.New("invalid authorization header format")
)

// github app
var ()

// project
var (
	ErrMissingNameInParams         = errors.New("name is missing in params")
	ErrMissingIdInParams           = errors.New("id is missing in params")
	ErrProjectNotFound             = errors.New("project not found")
	ErrProjectNameUniqueConstraint = errors.New("project with this name is already exist")
)

// deployment
var (
	ErrMissingTechstackNameInParams = errors.New("techstack name is missing in params")
)

// monitoring
var (
	ErrInvalidStartTimeFormat = errors.New("invalid start time format")
	ErrInvalidEndTimeFormat   = errors.New("invalid end time format")
)
