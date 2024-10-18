package service

import (
	"fmt"

	"github.com/NikolosHGW/goph-keeper/api/authpb"
	"github.com/NikolosHGW/goph-keeper/api/datapb"
	"github.com/NikolosHGW/goph-keeper/api/registerpb"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn           *grpc.ClientConn
	RegisterClient registerpb.RegisterClient
	AuthClient     authpb.AuthClient
	DataClient     datapb.DataServiceClient
}

func NewGRPCClient(serverAddress string, logger logger.CustomLogger) (*GRPCClient, error) {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.LogInfo("не удалось инициализировать клиент gRPC", err)

		return nil, fmt.Errorf("ошибка при инициализации gRPC клиента: %w", err)
	}

	registerClient := registerpb.NewRegisterClient(conn)
	authClient := authpb.NewAuthClient(conn)
	dataClient := datapb.NewDataServiceClient(conn)

	return &GRPCClient{
		conn:           conn,
		RegisterClient: registerClient,
		AuthClient:     authClient,
		DataClient:     dataClient,
	}, nil
}

func (c *GRPCClient) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("не удалось закрыть соединение клиента gRPC: %w", err)
	}

	return nil
}
