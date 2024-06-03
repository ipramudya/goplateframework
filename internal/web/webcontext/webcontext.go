package webcontext

import (
	"context"
	"errors"

	"github.com/goplateframework/internal/sdk/jsonwebtoken"
)

type ContextKey string

const (
	claimsKey         ContextKey = "claims_key"
	accountPayloadKey ContextKey = "account_payload_key"
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

func SetAccountPayload(ctx context.Context, p *jsonwebtoken.Payload) context.Context {
	return context.WithValue(ctx, accountPayloadKey, p)
}

func GetAccountPayload(ctx context.Context) (*jsonwebtoken.Payload, error) {
	val, ok := ctx.Value(accountPayloadKey).(*jsonwebtoken.Payload)

	if !ok {
		return &jsonwebtoken.Payload{}, errors.New("user not found")
	}

	return val, nil
}
