package service

import (
	"context"
	"errors"
	"testing"

	"github.com/NikolosHGW/goph-keeper/api/registerpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type MockRegisterClient struct {
	RegisterUserFunc func(
		ctx context.Context, req *registerpb.RegisterUserRequest, opts ...grpc.CallOption,
	) (*registerpb.RegisterUserResponse, error)
}

func (m *MockRegisterClient) RegisterUser(
	ctx context.Context, req *registerpb.RegisterUserRequest, opts ...grpc.CallOption,
) (*registerpb.RegisterUserResponse, error) {
	return m.RegisterUserFunc(ctx, req, opts...)
}

type mockLogger struct{}

func (n *mockLogger) LogInfo(message string, err error) {}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		mockRegister  func() *MockRegisterClient
		expectedToken string
		expectedErr   error
	}{
		{
			name:     "Успешная регистрация",
			login:    "user1",
			password: "pass1",
			mockRegister: func() *MockRegisterClient {
				return &MockRegisterClient{
					RegisterUserFunc: func(
						ctx context.Context, req *registerpb.RegisterUserRequest, opts ...grpc.CallOption,
					) (*registerpb.RegisterUserResponse, error) {
						return &registerpb.RegisterUserResponse{BearerToken: "token123"}, nil
					},
				}
			},
			expectedToken: "token123",
			expectedErr:   nil,
		},
		{
			name:     "Регистрация с ошибкой от RegisterClient",
			login:    "user2",
			password: "pass2",
			mockRegister: func() *MockRegisterClient {
				return &MockRegisterClient{
					RegisterUserFunc: func(
						ctx context.Context, req *registerpb.RegisterUserRequest, opts ...grpc.CallOption,
					) (*registerpb.RegisterUserResponse, error) {
						return nil, errors.New("grpc error")
					},
				}
			},
			expectedToken: "",
			expectedErr:   errors.New("ошибка при регистрации: grpc error"),
		},
		{
			name:     "Регистрация с пустым логином",
			login:    "",
			password: "pass3",
			mockRegister: func() *MockRegisterClient {
				return &MockRegisterClient{
					RegisterUserFunc: func(
						ctx context.Context, req *registerpb.RegisterUserRequest, opts ...grpc.CallOption,
					) (*registerpb.RegisterUserResponse, error) {
						return &registerpb.RegisterUserResponse{BearerToken: "token456"}, nil
					},
				}
			},
			expectedToken: "token456",
			expectedErr:   nil,
		},
		{
			name:     "Регистрация с пустым паролем",
			login:    "user4",
			password: "",
			mockRegister: func() *MockRegisterClient {
				return &MockRegisterClient{
					RegisterUserFunc: func(
						ctx context.Context, req *registerpb.RegisterUserRequest, opts ...grpc.CallOption,
					) (*registerpb.RegisterUserResponse, error) {
						return &registerpb.RegisterUserResponse{BearerToken: "token789"}, nil
					},
				}
			},
			expectedToken: "token789",
			expectedErr:   nil,
		},
		{
			name:     "Регистрация с контекстом, отмененным",
			login:    "user5",
			password: "pass5",
			mockRegister: func() *MockRegisterClient {
				return &MockRegisterClient{
					RegisterUserFunc: func(
						ctx context.Context, req *registerpb.RegisterUserRequest, opts ...grpc.CallOption,
					) (*registerpb.RegisterUserResponse, error) {
						return nil, context.Canceled
					},
				}
			},
			expectedToken: "",
			expectedErr:   context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGRPCClient := &GRPCClient{
				RegisterClient: tt.mockRegister(),
			}

			noOpLogger := &mockLogger{}

			authSvc := NewAuthService(mockGRPCClient, noOpLogger)

			token, err := authSvc.Register(context.Background(), tt.login, tt.password)

			assert.Equal(t, tt.expectedToken, token)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
