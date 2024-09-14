package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
)

const TokenExp = time.Hour * 5

func sendToken(w http.ResponseWriter, secretKey string, logger logger.CustomLogger, user *entity.User) {
	token, err := generateJWT(user, secretKey, logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(http.StatusOK)
}

func generateJWT(user *entity.User, secretKey string, logger logger.CustomLogger) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: user.ID,
	})

	if secretKey == "" {
		logger.LogInfo("для создании подписи токена секретный ключ пустой", fmt.Errorf("пустой secretKey"))
		return "", helper.ErrInternalServer
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.LogInfo("ошибки при создании подписи токена: ", err)
		return "", helper.ErrInternalServer
	}

	return tokenString, nil
}
