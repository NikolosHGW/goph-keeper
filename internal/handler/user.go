package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/request"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

var (
	ErrLoginAlreadyExists = errors.New("логин уже существует")
	ErrInvalidCredentials = errors.New("неверная пара логин/пароль")
)

type register interface {
	Register(context.Context, string, string) (*entity.User, error)
}

type UserHandler struct {
	registerUseCase register
	logger          logger.CustomLogger
}

func NewUserHandler(registerUseCase register, logger logger.CustomLogger) *UserHandler {
	return &UserHandler{
		registerUseCase: registerUseCase,
		logger:          logger,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	inputData, err := request.NewRegisterUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.registerUseCase.Register(r.Context(), inputData.Login, inputData.Password)
	if err != nil {
		if errors.Is(err, ErrLoginAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendToken(w, h, user)
}

func sendToken(w http.ResponseWriter, h *UserHandler, user *entity.User) {
	// token, err := h.userUseCase.GenerateJWT(user)
	// if err != nil {
	// 	http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
	// 	return
	// }
	fmt.Println(h, user)
	token := ""
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(http.StatusOK)
}
