package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
)

const maxPasswordLength = 72

// RegisterUser - реквест регистрации пользователя.
type RegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// NewRegisterUser - конструктор для валидации реквеста и преобразования в DTO.
func NewRegisterUser(r *http.Request, logger logger.CustomLogger) (*RegisterUser, error) {
	var registerData RegisterUser
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		logger.LogInfo("ошибка декодирования", err)
		return nil, helper.ErrInternalServer
	}

	err := validate(registerData)
	if err != nil {
		logger.LogInfo("ошибка валидации", err)
		return nil, fmt.Errorf("неверный формат запроса")
	}

	return &registerData, nil
}

func validate(registerData RegisterUser) error {
	if registerData.Login == "" || registerData.Password == "" {
		return errors.New("пустые логин и/или пароль")
	}
	if len([]byte(registerData.Password)) > maxPasswordLength {
		return errors.New("слишком длинный пароль")
	}

	return nil
}
