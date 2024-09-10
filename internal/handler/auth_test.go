package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/internal/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLogger struct{}

func (ml mockLogger) LogInfo(message string, err error) {}

type MockRegisterUseCase struct {
	mock.Mock
}

func (m *MockRegisterUseCase) Register(ctx context.Context, registerDTO *request.RegisterUser) (*entity.User, error) {
	args := m.Called(ctx, registerDTO)
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GenerateJWT(user *entity.User, secretKey string) (string, error) {
	args := m.Called(user, secretKey)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_RegisterUser_Success(t *testing.T) {
	mockRegisterUseCase := new(MockRegisterUseCase)
	mockUserService := new(MockUserService)
	mockLogger := &mockLogger{}
	secretKey := "mysecretkey"

	handler := NewAuthHandler(mockRegisterUseCase, mockUserService, mockLogger, secretKey)

	registerDTO := &request.RegisterUser{Login: "testuser", Password: "password123"}
	user := &entity.User{ID: 12345}

	mockRegisterUseCase.On("Register", mock.Anything, registerDTO).Return(user, nil)
	mockUserService.On("GenerateJWT", user, secretKey).Return("testtoken", nil)

	body := `{"login": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.RegisterUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Bearer testtoken", rr.Header().Get("Authorization"))
	mockRegisterUseCase.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestAuthHandler_RegisterUser_BadRequest(t *testing.T) {
	mockRegisterUseCase := new(MockRegisterUseCase)
	mockUserService := new(MockUserService)
	mockLogger := &mockLogger{}
	secretKey := "mysecretkey"

	handler := NewAuthHandler(mockRegisterUseCase, mockUserService, mockLogger, secretKey)

	body := `{"login": "", "password": ""}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.RegisterUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockRegisterUseCase.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestAuthHandler_RegisterUser_Conflict(t *testing.T) {
	mockRegisterUseCase := new(MockRegisterUseCase)
	mockUserService := new(MockUserService)
	mockLogger := &mockLogger{}
	secretKey := "mysecretkey"

	handler := NewAuthHandler(mockRegisterUseCase, mockUserService, mockLogger, secretKey)

	registerDTO := &request.RegisterUser{Login: "testuser", Password: "password123"}

	mockRegisterUseCase.
		On(
			"Register",
			mock.Anything,
			registerDTO,
		).
		Return(
			(*entity.User)(nil),
			helper.ErrLoginAlreadyExists,
		)

	body := `{"login": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.RegisterUser(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	mockRegisterUseCase.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestAuthHandler_RegisterUser_GenerateJWTError(t *testing.T) {
	mockRegisterUseCase := new(MockRegisterUseCase)
	mockUserService := new(MockUserService)
	mockLogger := &mockLogger{}
	secretKey := "mysecretkey"

	handler := NewAuthHandler(mockRegisterUseCase, mockUserService, mockLogger, secretKey)

	registerDTO := &request.RegisterUser{Login: "testuser", Password: "password123"}
	user := &entity.User{ID: 12345}

	mockRegisterUseCase.On("Register", mock.Anything, registerDTO).Return(user, nil)
	mockUserService.On("GenerateJWT", user, secretKey).Return("", errors.New("ошибка генерации токена"))

	body := `{"login": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.RegisterUser(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockRegisterUseCase.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestAuthHandler_RegisterUser_LongPassword(t *testing.T) {
	mockRegisterUseCase := new(MockRegisterUseCase)
	mockUserService := new(MockUserService)
	mockLogger := &mockLogger{}
	secretKey := "mysecretkey"

	handler := NewAuthHandler(mockRegisterUseCase, mockUserService, mockLogger, secretKey)

	longPassword := "aVeryLongPasswordThatExceedsTheLimitOf72BytesAndShouldTriggerAnErrorForValidation"

	body := `{"login": "testuser", "password": "` + longPassword + `"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.RegisterUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
