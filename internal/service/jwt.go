//nolint:exhaustruct
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken         = errors.New("ErrInvalidToken")
	ErrUnexpectedSignMethod = errors.New("ErrUnexpectedSignMethod")
	ErrInvalidTokenClaims   = errors.New("ErrInvalidTokenClaims")
	ErrInternalCodeGen      = errors.New("ErrInternalCodeGen")
)

type UserClaims struct {
	UUID string `json:"uuid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey []byte
	expiry    time.Duration
}

func NewJWTManager(secretKey string, expiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
		expiry:    expiry,
	}
}

func (m *JWTManager) GenerateToken(userUUID, role string) (string, error) {
	claims := UserClaims{
		UUID: userUUID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", ErrInternalCodeGen
	}

	return signed, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w (%v)", ErrUnexpectedSignMethod, token.Header["alg"])
			}

			return m.secretKey, nil
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("%w (%w)", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok || !token.Valid {
		return "", "", ErrInvalidTokenClaims
	}

	return claims.UUID, claims.Role, nil
}
