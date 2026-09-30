package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

// GitHub signs each delivery with HMAC-SHA256 over the raw request body and
// sends the hex digest in the X-Hub-Signature-256 header, prefixed "sha256=".
// The body must be the exact bytes received: any re-encoding before hashing
// invalidates the comparison, so callers must verify before decoding.
//
// https://docs.github.com/en/webhooks/using-webhooks/validating-webhook-deliveries
func VerifyGitHubSignature(secret string, body []byte, header string) error {
	if secret == "" || header == "" {
		return defined_error.ErrInvalidWebhookSignature
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// hmac.Equal is constant-time; a plain == comparison leaks via timing.
	if !hmac.Equal([]byte(expected), []byte(header)) {
		return defined_error.ErrInvalidWebhookSignature
	}

	return nil
}
