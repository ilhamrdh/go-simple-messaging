package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
)

type ClaimToken struct {
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	jwt.RegisteredClaims
}

var MapTokenType = map[string]time.Duration{
	"token":         time.Hour * 1,
	"refresh_token": time.Hour * 72,
}

func GenerateToken(username, fullname, tokenType string) (string, error) {
	secret := []byte(env.GetEnv("APP_SECRET_KEY", ""))

	claimToken := ClaimToken{
		Username: username,
		Fullname: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(MapTokenType[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)

	result, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return result, nil
}

func ValidateToken(token string) (*ClaimToken, error) {
	var ok bool
	claimToken := new(ClaimToken)
	secret := []byte(env.GetEnv("APP_SECRET_KEY", ""))

	jwtToken, err := jwt.ParseWithClaims(token, &ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to validate method jwt: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %v", err)
	}

	claimToken, ok = jwtToken.Claims.(*ClaimToken)

	if !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return claimToken, nil
}
