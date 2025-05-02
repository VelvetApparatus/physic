package token_service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"physk/internal/domain/aggregates/user"
	"time"
)

var (
	jwtSecretKey      = []byte("very-secret-key")
	InvalidTokenError = errors.New("invalid token")
)

type TokenService struct {
	secret []byte
}

func NewTokenService() TokenService {
	return TokenService{secret: jwtSecretKey}
}

func (ts *TokenService) NewToken(u user.User) (string, error) {
	payload := jwt.MapClaims{
		"id":         u.ID,
		"expires_at": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(ts.secret)
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}
	return t, nil
}

func (ts *TokenService) ValidateToken(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, InvalidTokenError
	}

	return jwtToken, nil
}
