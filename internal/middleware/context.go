package middleware

import "context"

type CtxKey int

const (
	UserIDKey = iota
	SessionIDKey
)

func createContext(userID uint, sessionID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, SessionIDKey, sessionID)
	return ctx
}
