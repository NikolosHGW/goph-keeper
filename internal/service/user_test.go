package service

import (
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/request"
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
