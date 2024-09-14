package request

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/helper"
)

func TestNewLogin(t *testing.T) {
	tests := []struct {
		name          string
		inputBody     string
		expectedError error
		expectedLogin *Login
	}{
		{
			name:          "Успешный запрос",
			inputBody:     `{"login":"user1","password":"password1"}`,
			expectedError: nil,
			expectedLogin: &Login{
				Login:    "user1",
				Password: "password1",
			},
		},
		{
			name:          "Некорректный JSON",
			inputBody:     `{"login":"user1","password":"password1"`,
			expectedError: helper.ErrInternalServer,
			expectedLogin: nil,
		},
		{
			name:          "Пустой логин",
			inputBody:     `{"login":"","password":"password1"}`,
			expectedError: errors.New("неверный формат запроса"),
			expectedLogin: nil,
		},
		{
			name:          "Пустой пароль",
			inputBody:     `{"login":"user1","password":""}`,
			expectedError: errors.New("неверный формат запроса"),
			expectedLogin: nil,
		},
		{
			name:          "Пустые логин и пароль",
			inputBody:     `{"login":"","password":""}`,
			expectedError: errors.New("неверный формат запроса"),
			expectedLogin: nil,
		},
		{
			name:          "Пустое тело запроса",
			inputBody:     ``,
			expectedError: helper.ErrInternalServer,
			expectedLogin: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(tc.inputBody)))
			req.Header.Set("Content-Type", "application/json")

			mockLogger := &mockLogger{}

			loginData, err := NewLogin(req, mockLogger)

			if tc.expectedError != nil {
				if err == nil {
					t.Errorf("ожидалась ошибка '%v', получили nil", tc.expectedError)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("ожидалась ошибка '%v', получили '%v'", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("ожидалась ошибка nil, получили '%v'", err)
				}
			}

			if tc.expectedLogin != nil {
				if loginData == nil {
					t.Errorf("ожидались данные логина, получили nil")
				} else {
					if loginData.Login != tc.expectedLogin.Login {
						t.Errorf("ожидался логин '%s', получили '%s'", tc.expectedLogin.Login, loginData.Login)
					}
					if loginData.Password != tc.expectedLogin.Password {
						t.Errorf("ожидался пароль '%s', получили '%s'", tc.expectedLogin.Password, loginData.Password)
					}
				}
			} else if loginData != nil {
				t.Errorf("ожидались данные логина nil, получили '%v'", loginData)
			}
		})
	}
}
