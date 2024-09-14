package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"golang.org/x/crypto/bcrypt"
)

type mockLoginRepo struct {
	FindByLoginFunc func(ctx context.Context, login string) (*entity.User, error)
}

func (m *mockLoginRepo) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	return m.FindByLoginFunc(ctx, login)
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name            string
		login           string
		password        string
		mockFindByLogin func(ctx context.Context, login string) (*entity.User, error)
		expectedUser    *entity.User
		expectedError   error
	}{
		{
			name:     "Успешная аутентификация",
			login:    "user1",
			password: "password1",
			mockFindByLogin: func(ctx context.Context, login string) (*entity.User, error) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
				return &entity.User{
					ID:       1,
					Login:    "user1",
					Password: string(hashedPassword),
				}, nil
			},
			expectedUser: &entity.User{
				ID:    1,
				Login: "user1",
			},
			expectedError: nil,
		},
		{
			name:     "Пользователь не найден",
			login:    "user1",
			password: "password1",
			mockFindByLogin: func(ctx context.Context, login string) (*entity.User, error) {
				return nil, helper.ErrInvalidCredentials
			},
			expectedUser:  nil,
			expectedError: helper.ErrInvalidCredentials,
		},
		{
			name:     "Неверный пароль",
			login:    "user1",
			password: "wrongpassword",
			mockFindByLogin: func(ctx context.Context, login string) (*entity.User, error) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
				return &entity.User{
					ID:       1,
					Login:    "user1",
					Password: string(hashedPassword),
				}, nil
			},
			expectedUser:  nil,
			expectedError: helper.ErrInvalidCredentials,
		},
		{
			name:     "Ошибка сервера при поиске пользователя",
			login:    "user1",
			password: "password1",
			mockFindByLogin: func(ctx context.Context, login string) (*entity.User, error) {
				return nil, errors.New("database error")
			},
			expectedUser:  nil,
			expectedError: errors.New("ошибка сервера: database error"),
		},
		{
			name:     "Некорректный хэш пароля",
			login:    "user1",
			password: "password1",
			mockFindByLogin: func(ctx context.Context, login string) (*entity.User, error) {
				return &entity.User{
					ID:       1,
					Login:    "user1",
					Password: "invalid-hash",
				}, nil
			},
			expectedUser:  nil,
			expectedError: helper.ErrInvalidCredentials,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockLoginRepo{
				FindByLoginFunc: tc.mockFindByLogin,
			}

			loginUsecase := NewLogin(mockRepo)

			user, err := loginUsecase.Authenticate(context.Background(), tc.login, tc.password)

			if tc.expectedError != nil {
				if err == nil {
					t.Errorf("ожидалась ошибка '%v', получили nil", tc.expectedError)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("ожидалась ошибка '%v', получили '%v'", tc.expectedError, err)
				}
				if user != nil {
					t.Errorf("ожидался пользователь nil, получили '%v'", user)
				}
			} else {
				if err != nil {
					t.Errorf("ожидалась ошибка nil, получили '%v'", err)
				}
				if user == nil {
					t.Errorf("ожидался пользователь, получили nil")
				} else {
					if user.ID != tc.expectedUser.ID {
						t.Errorf("ожидался ID пользователя '%d', получили '%d'", tc.expectedUser.ID, user.ID)
					}
					if user.Login != tc.expectedUser.Login {
						t.Errorf("ожидался логин '%s', получили '%s'", tc.expectedUser.Login, user.Login)
					}
				}
			}
		})
	}
}
