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
)

type mockLoginUsecase struct {
	AuthenticateFunc func(ctx context.Context, login string, password string) (*entity.User, error)
}

func (m *mockLoginUsecase) Authenticate(ctx context.Context, login string, password string) (*entity.User, error) {
	return m.AuthenticateFunc(ctx, login, password)
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name             string
		inputBody        string
		mockAuthenticate func(ctx context.Context, login string, password string) (*entity.User, error)
		expectedStatus   int
		expectedHeaders  map[string]string
		expectedError    string
		loginSecretKey   string
	}{
		{
			name:      "Успешный вход",
			inputBody: `{"login":"user1","password":"password1"}`,
			mockAuthenticate: func(ctx context.Context, login string, password string) (*entity.User, error) {
				return &entity.User{
					ID: 1,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Authorization": "Bearer ",
				"Content-Type":  "application/json",
			},
			loginSecretKey: "testsecretkey",
		},
		{
			name:      "Неверные учетные данные",
			inputBody: `{"login":"user1","password":"wrongpassword"}`,
			mockAuthenticate: func(ctx context.Context, login string, password string) (*entity.User, error) {
				return nil, helper.ErrInvalidCredentials
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  helper.ErrInvalidCredentials.Error(),
			loginSecretKey: "testsecretkey",
		},
		{
			name:      "Внутренняя ошибка сервера при аутентификации",
			inputBody: `{"login":"user1","password":"password1"}`,
			mockAuthenticate: func(ctx context.Context, login string, password string) (*entity.User, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database error",
			loginSecretKey: "testsecretkey",
		},
		{
			name:             "Некорректный запрос",
			inputBody:        `{"login":"user1","password":`,
			mockAuthenticate: nil,
			expectedStatus:   http.StatusBadRequest,
			loginSecretKey:   "testsecretkey",
		},
		{
			name:      "Пустой секретный ключ",
			inputBody: `{"login":"user1","password":"password1"}`,
			mockAuthenticate: func(ctx context.Context, login string, password string) (*entity.User, error) {
				return &entity.User{
					ID: 1,
				}, nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  helper.ErrInternalServer.Error(),
			loginSecretKey: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(tc.inputBody)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			mockUsecase := &mockLoginUsecase{}
			if tc.mockAuthenticate != nil {
				mockUsecase.AuthenticateFunc = tc.mockAuthenticate
			}

			mockLogger := &mockLogger{}

			loginHandler := NewLogin(mockUsecase, mockLogger, tc.loginSecretKey)

			loginHandler.Login(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("ожидался статус %d, получили %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedError != "" {
				respBody := rr.Body.String()
				if respBody != tc.expectedError+"\n" {
					t.Errorf("ожидалось сообщение об ошибке '%s', получили '%s'", tc.expectedError+"\n", respBody)
				}
			}

			for key, value := range tc.expectedHeaders {
				headerValue := rr.Header().Get(key)
				if headerValue == "" {
					t.Errorf("ожидался заголовок '%s'", key)
				}
				if value != "" && !containsPrefix(headerValue, value) {
					t.Errorf("ожидалось, что заголовок '%s' начнется с '%s', получили '%s'", key, value, headerValue)
				}
			}
		})
	}
}

func containsPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
