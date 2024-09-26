package service

import (
	"fmt"

	pb "github.com/NikolosHGW/goph-keeper/api/registerpb"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn           *grpc.ClientConn
	RegisterClient pb.RegisterClient
}

func NewGRPCClient(serverAddress string, logger logger.CustomLogger) (*GRPCClient, error) {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.LogInfo("не удалось инициализировать клиент gRPC", err)

		return nil, fmt.Errorf("ошибка при инициализации gRPC клиента: %w", err)
	}

	registerClient := pb.NewRegisterClient(conn)

	return &GRPCClient{
		conn:           conn,
		RegisterClient: registerClient,
	}, nil
}

func (c *GRPCClient) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("не удалось закрыть соединение клиента gRPC: %w", err)
	}

	return nil
}
