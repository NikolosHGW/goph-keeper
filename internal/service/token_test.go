package service

import (
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestToken_GenerateJWT_Success(t *testing.T) {
	mockLogger := &mockLogger{}
	secretKey := "supersecretkey"
	tokenService := NewToken(mockLogger, secretKey)

	user := &entity.User{
		ID: 1,
	}

	tokenString, err := tokenService.GenerateJWT(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	parsedToken, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)

	claims, ok := parsedToken.Claims.(*entity.Claims)
	assert.True(t, ok)
	assert.Equal(t, user.ID, claims.UserID)
}

func TestToken_GenerateJWT_EmptySecretKey(t *testing.T) {
	mockLogger := &mockLogger{}
	tokenService := NewToken(mockLogger, "")

	user := &entity.User{
		ID: 1,
	}

	tokenString, err := tokenService.GenerateJWT(user)

	assert.Empty(t, tokenString)
	assert.ErrorIs(t, err, helper.ErrInternalServer)
}
