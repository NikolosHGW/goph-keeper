package service

import (
	"fmt"
	"time"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/internal/request"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const TokenExp = time.Hour * 5

type User struct {
	logger logger.CustomLogger
}

func NewUser(logger logger.CustomLogger) *User {
	return &User{
		logger: logger,
	}
}

// GetUser - отдаёт сущность пользователя с захешированным паролем.
func (u *User) GetUser(registerDTO *request.RegisterUser) (*entity.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.LogInfo("ошибка при хэшировании пароля: ", err)
		return nil, helper.ErrInternalServer
	}

	user := &entity.User{
		Login:    registerDTO.Login,
		Password: string(passwordHash),
	}

	return user, nil
}

// GenerateJWT - генерирует токен.
func (u *User) GenerateJWT(user *entity.User, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: user.ID,
	})

	if secretKey == "" {
		u.logger.LogInfo("для создании подписи токена секретный ключ пустой", fmt.Errorf("пустой secretKey"))
		return "", helper.ErrInternalServer
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		u.logger.LogInfo("ошибки при создании подписи токена: ", err)
		return "", helper.ErrInternalServer
	}

	return tokenString, nil
}
