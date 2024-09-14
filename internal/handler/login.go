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

type loginUsecase interface {
	Authenticate(context.Context, string, string) (*entity.User, error)
}

type Login struct {
	loginUsecase loginUsecase
	logger       logger.CustomLogger
	secretKey    string
}

func NewLogin(loginUsecase loginUsecase, logger logger.CustomLogger, secretKey string) *Login {
	return &Login{
		loginUsecase: loginUsecase,
		logger:       logger,
		secretKey:    secretKey,
	}
}

func (l *Login) Login(w http.ResponseWriter, r *http.Request) {
	inputData, err := request.NewLogin(r, l.logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := l.loginUsecase.Authenticate(r.Context(), inputData.Login, inputData.Password)
	if err != nil {
		if errors.Is(err, helper.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendToken(w, l.secretKey, l.logger, user)
}
