package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewRegisterUser(r *http.Request) (*RegisterUser, error) {
	var registerData RegisterUser
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		return nil, fmt.Errorf("ошибка декодирования: %w", err)
	}

	err := validate(registerData)
	if err != nil {
		return nil, fmt.Errorf("ошибка валидации: %w", err)
	}

	return &registerData, nil
}

func validate(registerData RegisterUser) error {
	if registerData.Login == "" || registerData.Password == "" {
		return errors.New("неверный формат запроса")
	}

	return nil
}
