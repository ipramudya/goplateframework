package webcontext

import (
	"context"

	"github.com/goplateframework/internal/sdk/tokenutil"
)

type ContextKey string

const (
	claimsKey ContextKey = "claims_key"
	tokenKey  ContextKey = "token_key"
)

func SetClaims(ctx context.Context, cl *tokenutil.AccessTokenClaims) context.Context {
	return context.WithValue(ctx, claimsKey, cl)
}

func GetClaims(ctx context.Context) *tokenutil.AccessTokenClaims {
	val, ok := ctx.Value(claimsKey).(*tokenutil.AccessTokenClaims)

	if !ok {
		return &tokenutil.AccessTokenClaims{}
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
