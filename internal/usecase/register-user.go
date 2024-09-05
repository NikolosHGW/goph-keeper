package usecase

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
)

type RegisterUser struct {
	userService    interface{}
	userRepository interface{}
}

func (r *RegisterUser) Register(ctx context.Context, login, password string) (*entity.User, error) {
	fmt.Println(r.userRepository, r.userService, login, password)
	return &entity.User{}, nil
}
