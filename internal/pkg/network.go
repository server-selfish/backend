package pkg

import (
	"net/http"
	"strings"
)

func ReadIP(r *http.Request) string {
	xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	xri := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if xri != "" {
		return xri
	}
	host := strings.TrimSpace(r.RemoteAddr)
	if host == "" {
		return "unknown"
	}
	return host
}
