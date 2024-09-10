package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/internal/request"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

type register interface {
	Register(context.Context, *request.RegisterUser) (*entity.User, error)
}

type userServicer interface {
	GenerateJWT(*entity.User, string) (string, error)
}

type AuthHandler struct {
	registerUseCase register
	userService     userServicer
	logger          logger.CustomLogger
	secretKey       string
}

func NewAuthHandler(
	registerUseCase register,
	userService userServicer,
	logger logger.CustomLogger,
	secretKey string,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase,
		userService:     userService,
		logger:          logger,
		secretKey:       secretKey,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	registerDTO, err := request.NewRegisterUser(r, h.logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.registerUseCase.Register(r.Context(), registerDTO)
	if err != nil {
		if errors.Is(err, helper.ErrLoginAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendToken(w, h, user)
}

func sendToken(w http.ResponseWriter, h *AuthHandler, user *entity.User) {
	token, err := h.userService.GenerateJWT(user, h.secretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(http.StatusOK)
}
