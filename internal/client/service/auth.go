package service

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/goph-keeper/api/registerpb"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"google.golang.org/grpc"
)

type AuthService interface {
	Register(ctx context.Context, login, password string) (string, error)
}

type registerClient interface {
	RegisterUser(
		ctx context.Context,
		in *registerpb.RegisterUserRequest,
		opts ...grpc.CallOption,
	) (*registerpb.RegisterUserResponse, error)
}

type authService struct {
	registerClient registerClient
	logger         logger.CustomLogger
}

func NewAuthService(grpcClient *GRPCClient, logger logger.CustomLogger) *authService {
	return &authService{
		registerClient: grpcClient.RegisterClient,
		logger:         logger,
	}
}

func (s *authService) Register(ctx context.Context, login, password string) (string, error) {
	req := &registerpb.RegisterUserRequest{
		Login:    login,
		Password: password,
	}
	resp, err := s.registerClient.RegisterUser(ctx, req)
	if err != nil {
		s.logger.LogInfo("Ошибка регистрации", err)
		return "", fmt.Errorf("ошибка при регистрации: %w", err)
	}
	return resp.BearerToken, nil
}
