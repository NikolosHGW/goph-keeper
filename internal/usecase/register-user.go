package usecase

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/internal/request"
)

type userRepo interface {
	Save(context.Context, *entity.User) error
	ExistsByLogin(context.Context, string) (bool, error)
}

type userServicer interface {
	GetUser(*request.RegisterUser) (*entity.User, error)
}

type registerUser struct {
	userService userServicer
	userRepo    userRepo
}

// NewRegisterUser конструктор юзкейса регистрации пользователя.
func NewRegisterUser(userService userServicer, userRepo userRepo) *registerUser {
	return &registerUser{
		userService: userService,
		userRepo:    userRepo,
	}
}

// Register регистрация пользователя.
func (r *registerUser) Register(ctx context.Context, registerDTO *request.RegisterUser) (*entity.User, error) {
	isLoginExist, err := r.userRepo.ExistsByLogin(ctx, registerDTO.Login)
	if err != nil {
		return nil, fmt.Errorf("ошибка сервера: %w", err)
	}
	if isLoginExist {
		return nil, helper.ErrLoginAlreadyExists
	}

	user, err := r.userService.GetUser(registerDTO)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	if err := r.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("ошибка при сохранении пользователя: %w", err)
	}

	return user, nil
}
