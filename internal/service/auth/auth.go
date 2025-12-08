// Package auth provides interface with authentiaction methods.
package auth

import (
	"errors"
	"fmt"

	smodels "github.com/darrior/urlshortener/internal/service/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	ErrorTokenInvalid = errors.New("token is invalid")
	ErrorEmptyUserID  = errors.New("user_id is empty")
)

type Auth interface {
	NewToken() (tokenString string, err error)
	GetUserID(tokenString string) (userID string, err error)
	ValidateToken(tokenString string) (valid bool, err error)
}

type HS256Auth struct {
	key []byte
}

func NewHS256Auth(key string) Auth {
	return &HS256Auth{
		key: []byte(key),
	}
}

func (h *HS256Auth) NewToken() (string, error) {
	uuid, err := generateUUID()
	if err != nil {
		return "", err
	}

	claims := &smodels.Claims{
		UserID: uuid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(h.key)
	if err != nil {
		return "", fmt.Errorf("can not create signed string: %w", err)
	}

	return tokenString, nil
}

func (h *HS256Auth) GetUserID(tokenString string) (string, error) {
	claims := &smodels.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return h.key, nil
	})
	if err != nil {
		return "", fmt.Errorf("can not parse claims: %w", err)
	}

	if !token.Valid {
		return "", ErrorTokenInvalid
	}

	if claims.UserID == "" {
		return "", ErrorEmptyUserID
	}

	return claims.UserID, nil
}

func (h *HS256Auth) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (any, error) {
		return h.key, nil
	})
	if err != nil {
		return false, fmt.Errorf("cannot parse token")
	}

	return token.Valid, nil
}

func generateUUID() (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("can not generate UUID: %w", err)
	}

	return uuid.String(), nil
}
