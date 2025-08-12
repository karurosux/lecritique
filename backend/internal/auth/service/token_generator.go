package authservice

import (
	"crypto/rand"
	"encoding/hex"
	authinterface "kyooar/internal/auth/interface"
)

type TokenGenerator struct{}

func NewTokenGenerator() authinterface.TokenGenerator {
	return &TokenGenerator{}
}

func (t *TokenGenerator) GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (t *TokenGenerator) GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}