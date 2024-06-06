package webcontext

import (
	"context"

	"github.com/goplateframework/internal/sdk/tokenutil"
)

type contextKey string

const (
	accessClaimsKey  contextKey = "access_claims_key"
	refreshClaimsKey contextKey = "refresh_claims_key"
	accessKey        contextKey = "access_key"
	refreshKey       contextKey = "refresh_key"
)

func SetAccessTokenClaims(ctx context.Context, cl *tokenutil.AccessTokenClaims) context.Context {
	return context.WithValue(ctx, accessClaimsKey, cl)
}

func GetAccessTokenClaims(ctx context.Context) *tokenutil.AccessTokenClaims {
	val, ok := ctx.Value(accessClaimsKey).(*tokenutil.AccessTokenClaims)

	if !ok {
		return &tokenutil.AccessTokenClaims{}
	}

	return val
}

func SetAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, accessKey, token)
}

func GetAccessToken(ctx context.Context) string {
	val, ok := ctx.Value(accessKey).(string)

	if !ok {
		return ""
	}

	return val
}

func SetRefreshToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, refreshKey, token)
}

func GetRefreshToken(ctx context.Context) string {
	val, ok := ctx.Value(refreshKey).(string)

	if !ok {
		return ""
	}

	return val
}

func SetRefreshTokenClaims(ctx context.Context, cl *tokenutil.RefreshTokenClaims) context.Context {
	return context.WithValue(ctx, refreshClaimsKey, cl)
}

func GetRefreshTokenClaims(ctx context.Context) *tokenutil.RefreshTokenClaims {
	val, ok := ctx.Value(refreshClaimsKey).(*tokenutil.RefreshTokenClaims)

	if !ok {
		return &tokenutil.RefreshTokenClaims{}
	}

	return val
}
