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
)

type mockLogger struct{}

func (ml mockLogger) LogInfo(message string, err error) {}

type mockRegisterUseCase struct {
	RegisterFunc func(ctx context.Context, user *request.RegisterUser) (*entity.User, error)
}

func (m *mockRegisterUseCase) Register(ctx context.Context, user *request.RegisterUser) (*entity.User, error) {
	return m.RegisterFunc(ctx, user)
}

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            string
		mockRegister         func(ctx context.Context, user *request.RegisterUser) (*entity.User, error)
		expectedStatus       int
		expectedErrorMessage string
		expectedHeaders      map[string]string
		secretKey            string
	}{
		{
			name:      "Успешная регистрация",
			inputBody: `{"login":"user1","password":"password1"}`,
			mockRegister: func(ctx context.Context, user *request.RegisterUser) (*entity.User, error) {
				return &entity.User{
					ID:    1,
					Login: user.Login,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Authorization": "Bearer ",
				"Content-Type":  "application/json",
			},
			secretKey: "testsecretkey",
		},
		{
			name:      "Пользователь уже существует",
			inputBody: `{"login":"user1","password":"password1"}`,
			mockRegister: func(ctx context.Context, user *request.RegisterUser) (*entity.User, error) {
				return nil, helper.ErrLoginAlreadyExists
			},
			expectedStatus:       http.StatusConflict,
			expectedErrorMessage: helper.ErrLoginAlreadyExists.Error() + "\n",
			secretKey:            "testsecretkey",
		},
		{
			name:                 "Неверные данные запроса",
			inputBody:            `{"login":"user1","password":"password1"`,
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "внутренняя ошибка сервера\n",
			secretKey:            "testsecretkey",
		},
		{
			name:      "Внутренняя ошибка сервера при регистрации",
			inputBody: `{"login":"user1","password":"password1"}`,
			mockRegister: func(ctx context.Context, user *request.RegisterUser) (*entity.User, error) {
				return nil, errors.New("database error")
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedErrorMessage: "database error\n",
			secretKey:            "testsecretkey",
		},
		{
			name:      "Пустой секретный ключ",
			inputBody: `{"login":"user1","password":"password1"}`,
			mockRegister: func(ctx context.Context, user *request.RegisterUser) (*entity.User, error) {
				return &entity.User{
					ID:    1,
					Login: user.Login,
				}, nil
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedErrorMessage: helper.ErrInternalServer.Error() + "\n",
			secretKey:            "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(tc.inputBody)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			mockRegisterUseCase := &mockRegisterUseCase{}
			if tc.mockRegister != nil {
				mockRegisterUseCase.RegisterFunc = tc.mockRegister
			}

			mockLogger := &mockLogger{}

			handler := NewRegisterHandler(mockRegisterUseCase, mockLogger, tc.secretKey)

			handler.RegisterUser(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("ожидался статус %d, получили %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedErrorMessage != "" {
				respBody := rr.Body.String()
				if respBody != tc.expectedErrorMessage {
					t.Errorf("ожидалось сообщение об ошибке '%s', получили '%s'", tc.expectedErrorMessage, respBody)
				}
			}

			for key, expectedValue := range tc.expectedHeaders {
				actualValue := rr.Header().Get(key)
				if actualValue == "" {
					t.Errorf("ожидался заголовок '%s'", key)
				}
				if expectedValue != "" && !containsPrefix(actualValue, expectedValue) {
					t.Errorf(
						"ожидалось, что заголовок '%s' начнется с '%s', получили '%s'",
						key,
						expectedValue,
						actualValue,
					)
				}
			}
		})
	}
}
