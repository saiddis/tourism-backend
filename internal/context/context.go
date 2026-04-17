package context

import "context"

type contextKey string

const claimsKey contextKey = "claims"

type Claims struct {
	UserID int
	Email  string
	Role   string
}

func GetClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}
