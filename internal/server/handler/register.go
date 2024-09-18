package handler

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/NikolosHGW/goph-keeper/api/registerpb"
)

const maxPasswordLength = 72

type register interface {
	Handle(context.Context, *pb.RegisterUserRequest) (string, error)
}

// RegisterServer - структура gRPC сервера для регистрации пользователя.
type RegisterServer struct {
	pb.UnimplementedRegisterServer

	registerUseCase register
}

// NewRegisterServer - конструктор gRPC сервера для регистрации пользователя.
func NewRegisterServer(registerUseCase register) *RegisterServer {
	return &RegisterServer{registerUseCase: registerUseCase}
}

// RegisterUser - реализация RPC сервиса.
func (s *RegisterServer) RegisterUser(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {
	err := validateRegisterUserRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неправильный запрос: %v", err)
	}

	token, err := s.registerUseCase.Handle(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при регистрации пользователя: %v", err)
	}

	return &pb.RegisterUserResponse{
		BearerToken: token,
	}, nil
}

func validateRegisterUserRequest(req *pb.RegisterUserRequest) error {
	if req.Login == "" || req.Password == "" {
		return errors.New("пустые логин и/или пароль")
	}
	if len([]byte(req.Password)) > maxPasswordLength {
		return fmt.Errorf("пароль не может быть длиннее чем %d символов", maxPasswordLength)
	}

	return nil
}
