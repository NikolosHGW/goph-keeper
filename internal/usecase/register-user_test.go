package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/entity"
	"github.com/NikolosHGW/goph-keeper/internal/helper"
	"github.com/NikolosHGW/goph-keeper/internal/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

type MockUserServicer struct {
	mock.Mock
}

func (m *MockUserServicer) GetUser(registerDTO *request.RegisterUser) (*entity.User, error) {
	args := m.Called(registerDTO)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestRegisterUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockUserServicer := new(MockUserServicer)
	useCase := NewRegisterUser(mockUserServicer, mockUserRepo)

	registerDTO := &request.RegisterUser{
		Login:    "testuser",
		Password: "password123",
	}
	user := &entity.User{
		Login:    "testuser",
		Password: "hashedpassword",
		ID:       1,
	}

	mockUserRepo.On("ExistsByLogin", mock.Anything, registerDTO.Login).Return(false, nil)
	mockUserServicer.On("GetUser", registerDTO).Return(user, nil)
	mockUserRepo.On("Save", mock.Anything, user).Return(nil)

	result, err := useCase.Register(context.Background(), registerDTO)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockUserRepo.AssertExpectations(t)
	mockUserServicer.AssertExpectations(t)
}

func TestRegisterUser_LoginExists(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockUserServicer := new(MockUserServicer)
	useCase := NewRegisterUser(mockUserServicer, mockUserRepo)

	registerDTO := &request.RegisterUser{
		Login:    "testuser",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByLogin", mock.Anything, registerDTO.Login).Return(true, nil)

	result, err := useCase.Register(context.Background(), registerDTO)

	assert.ErrorIs(t, err, helper.ErrLoginAlreadyExists)
	assert.Nil(t, result)
	mockUserRepo.AssertExpectations(t)
}

func TestRegisterUser_ExistsByLoginError(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockUserServicer := new(MockUserServicer)
	useCase := NewRegisterUser(mockUserServicer, mockUserRepo)

	registerDTO := &request.RegisterUser{
		Login:    "testuser",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByLogin", mock.Anything, registerDTO.Login).Return(false, errors.New("database error"))

	result, err := useCase.Register(context.Background(), registerDTO)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ошибка сервера")
	mockUserRepo.AssertExpectations(t)
}

func TestRegisterUser_GetUserError(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockUserServicer := new(MockUserServicer)
	useCase := NewRegisterUser(mockUserServicer, mockUserRepo)

	registerDTO := &request.RegisterUser{
		Login:    "testuser",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByLogin", mock.Anything, registerDTO.Login).Return(false, nil)
	mockUserServicer.On("GetUser", registerDTO).Return((*entity.User)(nil), errors.New("failed to create user"))

	result, err := useCase.Register(context.Background(), registerDTO)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ошибка создания пользователя")
	mockUserRepo.AssertExpectations(t)
	mockUserServicer.AssertExpectations(t)
}

func TestRegisterUser_SaveError(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockUserServicer := new(MockUserServicer)
	useCase := NewRegisterUser(mockUserServicer, mockUserRepo)

	registerDTO := &request.RegisterUser{
		Login:    "testuser",
		Password: "password123",
	}
	user := &entity.User{
		Login:    "testuser",
		Password: "hashedpassword",
		ID:       1,
	}

	mockUserRepo.On("ExistsByLogin", mock.Anything, registerDTO.Login).Return(false, nil)
	mockUserServicer.On("GetUser", registerDTO).Return(user, nil)
	mockUserRepo.On("Save", mock.Anything, user).Return(errors.New("failed to save user"))

	result, err := useCase.Register(context.Background(), registerDTO)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ошибка при сохранении пользователя")
	mockUserRepo.AssertExpectations(t)
	mockUserServicer.AssertExpectations(t)
}
