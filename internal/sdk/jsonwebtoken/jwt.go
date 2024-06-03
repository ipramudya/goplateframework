package jsonwebtoken

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/goplateframework/config"
)

var (
	ExpiredTime     = time.Now().Add(time.Minute * 60) // 60 min from now
	Method          = jwt.GetSigningMethod(jwt.SigningMethodHS256.Name)
	ErrInvalidToken = errors.New("invalid token")
)

type Payload struct {
	Email     string `json:"email"`
	AccountID string `json:"account_id"`
	Role      string `json:"role"`
}

type Claims struct {
	jwt.RegisteredClaims
	Payload
}

func Generate(conf *config.Config, payload Payload) (string, error) {
	claims := &Claims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(ExpiredTime),
		},
	}
	token := jwt.NewWithClaims(Method, claims)
	tokenString, err := token.SignedString([]byte(conf.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Validate(conf *config.Config, authHeader string) (Claims, error) {
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, errors.New("expected authorization header format, \"Authorization: Bearer <token>\"")
	}

	var claims Claims

	token, err := jwt.ParseWithClaims(parts[1], &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, %v", token.Method)
		}
		return []byte(conf.Server.JwtSecretKey), nil
	})

	if err != nil {
		return Claims{}, fmt.Errorf("token signature is invalid")
	}

	if !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}
