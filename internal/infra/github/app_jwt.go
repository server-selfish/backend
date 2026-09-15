package github_infra

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
	"time"

	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

func generateAppJWTFromPEM(pkPEM, appId string) (string, error) {
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

func parseRSAPrivateKeyFromPEM(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, defined_error.ErrInvalidPrivateKey
	}

	if pk, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return pk, nil
	}

	privAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		if pk, ok := privAny.(*rsa.PrivateKey); ok {
			return pk, nil
		}
	}

	return nil, defined_error.ErrFailedParseRSAPrivateKey
}
