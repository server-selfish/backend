package schema

type AuthTokenPair struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresIn  int64  `json:"access_token_expires_in"`
	RefreshToken          string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in,omitempty"`
	TokenType             string `json:"token_type"`
}

type GithubCallbackResponse struct {
	Message string        `json:"message"`
	Data    AuthTokenPair `json:"data"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	Message string        `json:"message"`
	Data    AuthTokenPair `json:"data"`
}

type UserMeData struct {
	ID             string  `json:"id"`
	Provider       string  `json:"provider"`
	ProviderUserID int64   `json:"provider_user_id"`
	Username       string  `json:"username"`
	Email          *string `json:"email,omitempty"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
}

type MeResponse struct {
	Message string     `json:"message"`
	Data    UserMeData `json:"data"`
}

type GithubUserProfile struct {
	ID        int64   `json:"id"`
	Login     string  `json:"login"`
	Email     *string `json:"email"`
	AvatarURL *string `json:"avatar_url"`
}

type GithubOAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type AuthErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
