package repository

import (
	"context"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type User struct {
	db     *sqlx.DB
	logger logger.CustomLogger
}

func NewUser(db *sqlx.DB, logger logger.CustomLogger) *User {
	return &User{db: db, logger: logger}
}

func (r *User) Save(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowxContext(ctx, query, user.Login, user.Password).Scan(&user.ID)
	if err != nil {
		r.logger.LogInfo("ошибка при сохранении пользователя", err)
		return helper.ErrInternalServer
	}

	return nil
}

func (r *User) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login=$1)`
	err := r.db.QueryRowxContext(ctx, query, login).Scan(&exists)
	if err != nil {
		r.logger.LogInfo("не получилось записать результат запроса в переменную", err)
		return false, helper.ErrInternalServer
	}
	return exists, nil
}
