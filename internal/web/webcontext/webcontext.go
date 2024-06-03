package webcontext

import (
	"context"

	"github.com/goplateframework/internal/sdk/jsonwebtoken"
)

type ContextKey string

const (
	claimsKey ContextKey = "claims_key"
	tokenKey  ContextKey = "token_key"
)

func SetClaims(ctx context.Context, cl *jsonwebtoken.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, cl)
}

func GetClaims(ctx context.Context) *jsonwebtoken.Claims {
	val, ok := ctx.Value(claimsKey).(*jsonwebtoken.Claims)

	if !ok {
		return &jsonwebtoken.Claims{}
	}

	return val
}

func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func GetToken(ctx context.Context) string {
	val, ok := ctx.Value(tokenKey).(string)

	if !ok {
		return ""
	}

	return val
}
