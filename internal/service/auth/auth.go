// Package auth provides interface with authentiaction methods.
package auth

import (
	"errors"
	"fmt"

	smodels "github.com/darrior/urlshortener/internal/service/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Auth interface {
	SignClaims(claims *smodels.Claims) (tokenString string, err error)
	ValidateToken(tokenString string) (claims *smodels.Claims, err error)
	NewUserID() (userID string, err error)
}

type HS256Auth struct {
	key []byte
}

func NewHS256Auth(key string) Auth {
	return &HS256Auth{
		key: []byte(key),
	}
}

func (h *HS256Auth) SignClaims(claims *smodels.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(h.key)
	if err != nil {
		return "", fmt.Errorf("can not create signed string: %w", err)
	}

	return tokenString, nil
}

func (h *HS256Auth) ValidateToken(tokenString string) (*smodels.Claims, error) {
	claims := &smodels.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return h.key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("can not parse claims: %w", err)
	}

	if !token.Valid {
		errorTokenInvalid := errors.New("token is invalid")
		return nil, errorTokenInvalid
	}

	return claims, nil
}

func (h *HS256Auth) NewUserID() (string, error) {
	userID, err := generateUUID()
	if err != nil {
		return "", fmt.Errorf("can not generate uuid: %w", err)
	}

	return userID, err
}

func generateUUID() (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("can not generate UUID: %w", err)
	}

	return uuid.String(), nil
}
