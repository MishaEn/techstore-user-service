package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager struct {
	privateKey []byte
	publicKey  []byte
}

func NewTokenManager(privateKey, publicKey string) *TokenManager {
	return &TokenManager{
		privateKey: []byte(privateKey),
		publicKey:  []byte(publicKey),
	}
}

func (m *TokenManager) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Токен живет 15 минут
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.parsePrivateKey())
}

func (m *TokenManager) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.parsePublicKey(), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid subject")
	}

	return uuid.Parse(sub)
}

// Вспомогательные функции для парсинга ключей
func (m *TokenManager) parsePrivateKey() interface{} {
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(m.privateKey)
	return key
}

func (m *TokenManager) parsePublicKey() interface{} {
	key, _ := jwt.ParseRSAPublicKeyFromPEM(m.publicKey)
	return key
}
