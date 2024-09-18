package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/NikolosHGW/goph-keeper/api/registerpb"
	"github.com/NikolosHGW/goph-keeper/internal/client/infrastructure/config"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := config.NewConfig()

	myLogger, err := logger.NewLogger("info")
	if err != nil {
		log.Fatalf("ошибка инициализации логгер: %v", err)
	}

	conn, err := grpc.NewClient(config.GetServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		myLogger.LogInfo("не удалось инициализировать клиент gRPC", err)
		log.Fatal(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			myLogger.LogInfo("не удалось закрыть соединение клиента gRPC", err)
		}
	}()

	c := pb.NewRegisterClient(conn)
	resp, err := c.RegisterUser(context.Background(), &pb.RegisterUserRequest{Login: "ololo", Password: "qq!"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
