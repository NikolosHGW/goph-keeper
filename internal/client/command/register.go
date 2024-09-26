package command

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/goph-keeper/internal/client/entity"
	"github.com/NikolosHGW/goph-keeper/internal/client/service"
)

type RegisterCommand struct {
	authService service.AuthService
	tokenHolder *entity.TokenHolder
}

func NewRegisterCommand(authService service.AuthService, tokenHolder *entity.TokenHolder) *RegisterCommand {
	return &RegisterCommand{authService: authService, tokenHolder: tokenHolder}
}

func (c *RegisterCommand) Name() string {
	return "register"
}

func (c *RegisterCommand) Execute() error {
	var login, password string
	fmt.Print("Введите login: ")
	_, err := fmt.Scanln(&login)
	if err != nil {
		return fmt.Errorf("ошибка ввода логина: %w", err)
	}
	fmt.Print("Введите password: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		return fmt.Errorf("ошибка ввода пароля: %w", err)
	}

	token, err := c.authService.Register(context.Background(), login, password)
	if err != nil {
		return fmt.Errorf("ошибка решистрации: %w", err)
	}

	c.tokenHolder.Token = token
	fmt.Println("Регистрация прошла успешно.")
	return nil
}
