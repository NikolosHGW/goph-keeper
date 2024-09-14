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

type RegisterHandler struct {
	registerUseCase register
	logger          logger.CustomLogger
	secretKey       string
}

func NewRegisterHandler(
	registerUseCase register,
	logger logger.CustomLogger,
	secretKey string,
) *RegisterHandler {
	return &RegisterHandler{
		registerUseCase: registerUseCase,
		logger:          logger,
		secretKey:       secretKey,
	}
}

func (h *RegisterHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
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

	sendToken(w, h.secretKey, h.logger, user)
}
