package token_service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	jwtSecretKey      = []byte("very-secret-key")
	InvalidTokenError = errors.New("invalid token")
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret []byte
}

func NewTokenService() TokenService {
	return TokenService{secret: jwtSecretKey}
}

func (ts *TokenService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(ts.secret)
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}
	return t, nil
}

func (ts *TokenService) ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return ts.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, InvalidTokenError
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, InvalidTokenError

}
