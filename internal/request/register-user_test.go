package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (ml mockLogger) LogInfo(message string, err error) {}

func TestNewRegisterUser_Success(t *testing.T) {
	mockLogger := new(mockLogger)

	body := `{"login": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	registerUser, err := NewRegisterUser(req, mockLogger)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", registerUser.Login)
	assert.Equal(t, "password123", registerUser.Password)
}

func TestNewRegisterUser_DecodeError(t *testing.T) {
	mockLogger := new(mockLogger)

	body := `{"login": "testuser", "password": }`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	registerUser, err := NewRegisterUser(req, mockLogger)

	assert.Nil(t, registerUser)
	assert.Error(t, err)
	assert.Equal(t, helper.ErrInternalServer, err)
}

func TestNewRegisterUser_EmptyLoginOrPassword(t *testing.T) {
	mockLogger := new(mockLogger)

	body := `{"login": "", "password": ""}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	registerUser, err := NewRegisterUser(req, mockLogger)

	assert.Nil(t, registerUser)
	assert.Error(t, err)
	assert.Equal(t, "неверный формат запроса", err.Error())
}

func TestNewRegisterUser_LongPassword(t *testing.T) {
	mockLogger := new(mockLogger)

	longPassword := "aVeryLongPasswordThatExceedsTheLimitOf72BytesAndShouldTriggerAnErrorForValidation"
	body := `{"login": "testuser", "password": "` + longPassword + `"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	registerUser, err := NewRegisterUser(req, mockLogger)

	assert.Nil(t, registerUser)
	assert.Error(t, err)
	assert.Equal(t, "неверный формат запроса", err.Error())
}
