package tokenutil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/goplateframework/config"
)

type TokenStrategy int

const (
	AccessTokenStrategy TokenStrategy = iota
	RefreshTokenStrategy
)

var (
	AccessTokenExpiredTime  = time.Now().Add(time.Minute * 10) // 60 min from now
	RefreshTokenExpiredTime = time.Now().AddDate(0, 0, 30)     // 30 days from now
	Method                  = jwt.GetSigningMethod(jwt.SigningMethodHS256.Name)
	ErrInvalidToken         = errors.New("invalid token")
)

type AccessTokenPayload struct {
	Email     string `json:"email"`
	Role      string `json:"role"`
	AccountID string `json:"account_id"`
}

type RefreshTokenPayload struct {
	AccountID string `json:"account_id"`
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	AccessTokenPayload
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	RefreshTokenPayload
}

func GenerateAccess(conf *config.Config, payload AccessTokenPayload) (string, error) {
	claims := &AccessTokenClaims{
		AccessTokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(AccessTokenExpiredTime),
		},
	}

	token := jwt.NewWithClaims(Method, claims)
	tokenString, err := token.SignedString([]byte(conf.Server.JWTAccessTokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefresh(conf *config.Config, payload RefreshTokenPayload) (string, error) {
	claims := &RefreshTokenClaims{
		RefreshTokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(RefreshTokenExpiredTime),
		},
	}

	token := jwt.NewWithClaims(Method, claims)
	tokenString, err := token.SignedString([]byte(conf.Server.JWTRefreshTokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateAccess(conf *config.Config, requestToken string) (*AccessTokenClaims, error) {
	claims := new(AccessTokenClaims)

	token, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, %v", token.Method)
		}
		return []byte(conf.Server.JWTAccessTokenSecret), nil
	})

	if err != nil {
		return &AccessTokenClaims{}, fmt.Errorf("token signature is invalid")
	}

	if !token.Valid {
		return &AccessTokenClaims{}, ErrInvalidToken
	}

	return claims, nil
}

func ValidateRefresh(conf *config.Config, requestToken string) error {
	claims := new(RefreshTokenClaims)

	token, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, %v", token.Method)
		}
		return []byte(conf.Server.JWTRefreshTokenSecret), nil
	})

	if err != nil {
		return fmt.Errorf("token signature is invalid")
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

func RemainingTime(claims *jwt.RegisteredClaims) time.Duration {
	t := claims.ExpiresAt.Time
	return time.Until(t)
}

func ExtractBearerToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("expected authorization header format, \"Authorization: Bearer <token>\"")
	}

	return parts[1], nil
}
