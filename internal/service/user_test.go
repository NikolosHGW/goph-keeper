package service

import (
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/internal/request"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockLogger struct{}

func (ml mockLogger) LogInfo(message string, err error) {}

func TestUser_GetUser_Success(t *testing.T) {
	mockLogger := &mockLogger{}
	userService := NewUser(mockLogger)

	registerDTO := &request.RegisterUser{
		Login:    "testuser",
		Password: "password123",
	}

	user, err := userService.GetUser(registerDTO)

	assert.NoError(t, err)
	assert.Equal(t, registerDTO.Login, user.Login)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(registerDTO.Password))
	assert.NoError(t, err, "Пароль должен быть корректно захеширован")
}

func TestUser_GetUser_HashError(t *testing.T) {
	mockLogger := &mockLogger{}
	userService := NewUser(mockLogger)

	registerDTO := &request.RegisterUser{
		Login:    "testuser",
		Password: string(make([]byte, 100)),
	}

	user, err := userService.GetUser(registerDTO)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUser_GenerateJWT_Success(t *testing.T) {
	mockLogger := &mockLogger{}
	userService := NewUser(mockLogger)

	user := &entity.User{
		ID: 12345,
	}
	secretKey := "mysecretkey"

	token, err := userService.GenerateJWT(user, secretKey)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.ParseWithClaims(token, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*entity.Claims)
	assert.True(t, ok)
	assert.Equal(t, user.ID, claims.UserID)
}

func TestUser_GenerateJWT_EmptySecretKey(t *testing.T) {
	mockLogger := &mockLogger{}
	userService := NewUser(mockLogger)

	user := &entity.User{
		ID: 12345,
	}
	secretKey := ""

	token, err := userService.GenerateJWT(user, secretKey)

	assert.Error(t, err)
	assert.Equal(t, helper.ErrInternalServer, err)
	assert.Empty(t, token)
}
