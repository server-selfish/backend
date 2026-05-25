package pkg

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
	"github.com/spf13/viper"
)

type TokenManager struct {
	secret        []byte
	issuer        string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	refreshLength int
}

type AccessTokenClaims struct {
	UserID    string `json:"uid"`
	SessionID string `json:"sid"`
	Provider  string `json:"provider"`
	jwt.RegisteredClaims
}

func NewTokenManager(logger zerolog.Logger) (*TokenManager, error) {
	jwtSecret := viper.GetString("auth.jwt.secret")
	jwtIssuer := viper.GetString("auth.jwt.issuer")
	accessMinutes := viper.GetInt("auth.jwt.access_token_ttl_minutes")
	refreshDays := viper.GetInt("auth.jwt.refresh_token_ttl_days")

	if accessMinutes <= 0 {
		accessMinutes = 15
	}
	if refreshDays <= 0 {
		refreshDays = 7
	}

	if jwtSecret == "" {
		logger.Fatal().Msg("jwt secret is required")
	}
	if jwtIssuer == "" {
		jwtIssuer = "selfish-backend"
	}

	return &TokenManager{
		secret:        []byte(jwtSecret),
		issuer:        jwtIssuer,
		accessTTL:     time.Duration(accessMinutes) * time.Minute,
		refreshTTL:    time.Duration(refreshDays) * 24 * time.Hour,
		refreshLength: 64,
	}, nil
}

func (tm *TokenManager) GenerateAccessToken(userID, sessionID, provider string) (string, time.Time, error) {
	if userID == "" {
		return "", time.Time{}, defined_error.ErrMissingUserId
	}
	if sessionID == "" {
		return "", time.Time{}, defined_error.ErrMissingSessionId
	}
	if provider == "" {
		return "", time.Time{}, defined_error.ErrMissingProvider
	}

	now := time.Now().UTC()
	exp := now.Add(tm.accessTTL)

	claims := AccessTokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		Provider:  provider,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(tm.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%s-%w", defined_error.ErrSignAccessToken.Error(), err)
	}

	return signed, exp, nil
}

func (tm *TokenManager) ParseAccessToken(tokenString string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method == nil || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("%s: %v", defined_error.ErrUnexpectedSigningMethod.Error(), t.Header["alg"])
		}
		return tm.secret, nil
	}, jwt.WithIssuer(tm.issuer), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", defined_error.ErrParseAccessToken.Error(), err)
	}
	if !token.Valid {
		return nil, defined_error.ErrInvalidAccesstoken
	}

	if claims.UserID == "" || claims.SessionID == "" || claims.Provider == "" {
		return nil, defined_error.ErrInvaliTokenPayload
	}

	return claims, nil
}

func (tm *TokenManager) NewRefreshToken() (raw string, tokenHash string, expiresAt time.Time, err error) {
	buf := make([]byte, tm.refreshLength)
	if _, err = rand.Read(buf); err != nil {
		return "", "", time.Time{}, fmt.Errorf("%s: %w", defined_error.ErrGenerateRefreshToken.Error(), err)
	}

	// URL-safe opaque token
	raw = base64.RawURLEncoding.EncodeToString(buf)
	tokenHash = HashToken(raw)
	expiresAt = time.Now().UTC().Add(tm.refreshTTL)
	return raw, tokenHash, expiresAt, nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func GenerateAppJWTFromPEM(pkPEM, appId string) (string, error) {
	keyPEM := strings.TrimSpace(pkPEM)
	if keyPEM == "" {
		return "", defined_error.ErrMissingPrivateKey
	}
	privateKey, err := parseRSAPrivateKeyFromPEM(keyPEM)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	header := `{"alg":"RS256","typ":"JWT"}`

	appIDInt, err := strconv.ParseInt(appId, 10, 64)
	if err != nil {
		return "", defined_error.ErrInvalidGithubAppID
	}

	payloadMap := map[string]any{
		"iat": now.Add(-30 * time.Second).Unix(),
		"exp": now.Add(9 * time.Minute).Unix(),
		"iss": appIDInt,
	}
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		return "", fmt.Errorf("%s: %w", defined_error.ErrMarshalError.Error(), err)
	}

	unsigned := base64.RawURLEncoding.EncodeToString([]byte(header)) + "." + base64.RawURLEncoding.EncodeToString(payloadBytes)
	hashed := sha256.Sum256([]byte(unsigned))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("%s: %w", defined_error.ErrSignTokenError.Error(), err)
	}

	return unsigned + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}
