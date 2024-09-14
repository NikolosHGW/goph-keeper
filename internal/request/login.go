package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
)

// Login - реквест авторизации пользователя.
type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// NewLogin - конструктор для валидации реквеста и преобразования в DTO.
func NewLogin(r *http.Request, logger logger.CustomLogger) (*Login, error) {
	var loginData Login
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		logger.LogInfo("ошибка декодирования", err)
		return nil, helper.ErrInternalServer
	}

	err := validateLogin(loginData)
	if err != nil {
		logger.LogInfo("ошибка валидации", err)
		return nil, fmt.Errorf("неверный формат запроса")
	}

	return &loginData, nil
}

func validateLogin(loginData Login) error {
	if loginData.Login == "" || loginData.Password == "" {
		return errors.New("пустые логин и/или пароль")
	}

	return nil
}
