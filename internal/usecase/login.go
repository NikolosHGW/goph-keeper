package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"golang.org/x/crypto/bcrypt"
)

type loginRepo interface {
	FindByLogin(context.Context, string) (*entity.User, error)
}

type login struct {
	loginRepo loginRepo
}

func NewLogin(loginRepo loginRepo) *login {
	return &login{
		loginRepo: loginRepo,
	}
}

func (s *login) Authenticate(ctx context.Context, login, password string) (*entity.User, error) {
	user, err := s.loginRepo.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, helper.ErrInvalidCredentials) {
			return nil, helper.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("ошибка сервера: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, helper.ErrInvalidCredentials
	}

	return user, nil
}
