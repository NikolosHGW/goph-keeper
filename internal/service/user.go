package service

import (
	"errors"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/request"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	logger logger.CustomLogger
}

func NewUser(logger logger.CustomLogger) *User {
	return &User{
		logger: logger,
	}
}

func (u *User) GetUser(registerDTO *request.RegisterUser) (*entity.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.LogInfo("ошибка при хэшировании пароля: ", err)
		return nil, errors.New("временная ошибка сервиса, попробуйте ещё раз позже")
	}

	user := &entity.User{
		Login:    registerDTO.Login,
		Password: string(passwordHash),
	}

	return user, nil
}
