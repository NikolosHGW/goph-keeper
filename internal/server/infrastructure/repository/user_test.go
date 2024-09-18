package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NikolosHGW/goph-keeper/internal/server/entity"
	"github.com/NikolosHGW/goph-keeper/internal/server/helper"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (l *mockLogger) LogInfo(message string, err error) {}

func TestUser_Save_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mockLogger := new(mockLogger)
	repo := NewUser(sqlxDB, mockLogger)

	user := &entity.User{
		Login:    "testuser",
		Password: "password123",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", "password123").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.Save(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, int(1), user.ID)
}

func TestUser_Save_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mockLogger := new(mockLogger)
	repo := NewUser(sqlxDB, mockLogger)

	user := &entity.User{
		Login:    "testuser",
		Password: "password123",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", "password123").
		WillReturnError(errors.New("some error"))

	err = repo.Save(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, helper.ErrInternalServer, err)
}

func TestUser_ExistsByLogin_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mockLogger := new(mockLogger)
	repo := NewUser(sqlxDB, mockLogger)

	login := "testuser"

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(login).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsByLogin(context.Background(), login)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUser_ExistsByLogin_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mockLogger := new(mockLogger)
	repo := NewUser(sqlxDB, mockLogger)

	login := "testuser"

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(login).
		WillReturnError(errors.New("some error"))

	exists, err := repo.ExistsByLogin(context.Background(), login)

	assert.Error(t, err)
	assert.False(t, exists)
	assert.Equal(t, helper.ErrInternalServer, err)
}
