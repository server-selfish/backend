package pkg

import "context"

type contextKey string

const (
	ContextKeyUserID    contextKey = "auth_user_id"
	ContextKeySessionID contextKey = "auth_session_id"
	ContextKeyProvider  contextKey = "auth_provider"
)

func WithAuthUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

func WithAuthSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, ContextKeySessionID, sessionID)
}

func WithAuthProvider(ctx context.Context, provider string) context.Context {
	return context.WithValue(ctx, ContextKeyProvider, provider)
}

func AuthUserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}

func AuthSessionIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ContextKeySessionID).(string)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}

func AuthProviderFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ContextKeyProvider).(string)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}
