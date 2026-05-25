package defined_error

import "errors"

// common
var (
	ErrStringUUIDTypeCasting = errors.New("string-uuid type casting error")
	ErrStringIntTypeCasting  = errors.New("string-int type casting error")
)

// auth internal error
var (
	ErrMissingCodeInQueryParams       = errors.New("code params is missing in github callback")
	ErrMissingUserIdInContext         = errors.New("user id is empty from context")
	ErrExecuteTemplateError           = errors.New("github template execute error")
	ErrInvalidGithubProfile           = errors.New("invalid github profile response")
	ErrUpsertGithubUser               = errors.New("upsert github user error")
	ErrInvalidRefreshTokenInSession   = errors.New("invalid refresh token, not in the session")
	ErrRotateRefreshToken             = errors.New("rotate refresh token error")
	ErrGenerateAccessToken            = errors.New("generate access token error")
	ErrGenerateRefreshToken           = errors.New("generate refresh token error")
	ErrGetAuthSessionByRefreshToken   = errors.New("get auth by refresh token error")
	ErrRevokeSession                  = errors.New("revoke session error")
	ErrGetUserData                    = errors.New("get user data error")
	ErrEmptyAccessToken               = errors.New("empty access token")
	ErrInvalidOrExpireAccessToken     = errors.New("invalid or expired access token")
	ErrCreateAuthSession              = errors.New("create auth session error")
	ErrGithubAccessTokenRequestFailed = errors.New("github token request failed")
	ErrGithubTokenExchangeFailed      = errors.New("github token exchange failed")
	ErrDecodeGithubTokenResponse      = errors.New("decode github token failed")
	ErrEmptyGithubAccessToken         = errors.New("empty github access token")
	ErrGithubProfileRequestFailed     = errors.New("github profile request failed")
	ErrDecodeGithubProfileResponse    = errors.New("decode github profile failed")
)

// github app internal error
var (
	ErrMissingInstallationIDOrStateParams     = errors.New("missing installation id and/or state in params")
	ErrMissingAppIDOrSlug                     = errors.New("missing app_id and/or app_slug")
	ErrMissingGithubAppPrivateKey             = errors.New("missing github app private key")
	ErrMissingGithubAppCallbackURI            = errors.New("missing github app callbackURI")
	ErrParseUrl                               = errors.New("parse url error")
	ErrInvalidOrExpiredState                  = errors.New("invalid or expired state")
	ErrUpsertGithubInstallation               = errors.New("upsert github installation error")
	ErrGetListGithubInstallation              = errors.New("get list of github installation error")
	ErrGetGithubInstallationNotFound          = errors.New("installation not found")
	ErrGetGithubInstallation                  = errors.New("query installation error")
	ErrListInstallationRepositoriesRequest    = errors.New("list installation repositories request failed")
	ErrDecodeInstallationRepositoriesResponse = errors.New("decode installation repositories response error")
	ErrSetCacheGeneric                        = errors.New("set data to key error")
	// ErrGenerateStateCode                  = errors.New("generate state ")
)

// token internal error
var (
	ErrMissingUserId           = errors.New("missing user id")
	ErrMissingSessionId        = errors.New("missing session id")
	ErrMissingProvider         = errors.New("missing provider")
	ErrSignAccessToken         = errors.New("error sign access token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrParseAccessToken        = errors.New("error parse access token")
	ErrInvalidAccesstoken      = errors.New("invalid access token")
	ErrInvaliTokenPayload      = errors.New("invalid token payload")
	ErrMissingPrivateKey       = errors.New("missing private key")
	ErrInvalidGithubAppID      = errors.New("invalid github app id")
	ErrMarshalError            = errors.New("marshal error")
	ErrSignTokenError          = errors.New("app jwt sign error")
)

// deployment
var (
	ErrMissinProjectId          = errors.New("missing project id")
	ErrActiveDeploymentNotFound = errors.New("active deployment is not found")
	ErrDeploymentNotFound       = errors.New("deployment is not found")
)

// container
var (
	ErrContainerNotFound = errors.New("container is not found")
)

// parser internal error
var (
	ErrInvalidPrivateKey        = errors.New("invalid private key")
	ErrFailedParseRSAPrivateKey = errors.New("failed to parse rsa private key")
)
