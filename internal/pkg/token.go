package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func NewTokenManager(secret, issuer string, accessTTL, refreshTTL time.Duration) (*TokenManager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if issuer == "" {
		issuer = "selfish-backend"
	}
	if accessTTL <= 0 {
		return nil, errors.New("access token ttl must be greater than zero")
	}
	if refreshTTL <= 0 {
		return nil, errors.New("refresh token ttl must be greater than zero")
	}

	return &TokenManager{
		secret:        []byte(secret),
		issuer:        issuer,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		refreshLength: 64,
	}, nil
}

func (tm *TokenManager) GenerateAccessToken(userID, sessionID, provider string) (string, time.Time, error) {
	if userID == "" {
		return "", time.Time{}, errors.New("user id is required")
	}
	if sessionID == "" {
		return "", time.Time{}, errors.New("session id is required")
	}
	if provider == "" {
		return "", time.Time{}, errors.New("provider is required")
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
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	return signed, exp, nil
}

func (tm *TokenManager) ParseAccessToken(tokenString string) (*AccessTokenClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is required")
	}

	claims := &AccessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method == nil || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return tm.secret, nil
	}, jwt.WithIssuer(tm.issuer), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("parse access token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	if claims.UserID == "" || claims.SessionID == "" || claims.Provider == "" {
		return nil, errors.New("invalid token payload")
	}

	return claims, nil
}

func (tm *TokenManager) NewRefreshToken() (raw string, tokenHash string, expiresAt time.Time, err error) {
	buf := make([]byte, tm.refreshLength)
	if _, err = rand.Read(buf); err != nil {
		return "", "", time.Time{}, fmt.Errorf("generate refresh token: %w", err)
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
