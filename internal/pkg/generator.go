package pkg

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

func GenerateState(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("state length must be positive")
	}

	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	// URL-safe, no padding; good for query param usage.
	return base64.RawURLEncoding.EncodeToString(b), nil
}
