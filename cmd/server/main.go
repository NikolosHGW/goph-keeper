package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikolosHGW/goph-keeper/api/registerpb"
	"github.com/NikolosHGW/goph-keeper/internal/server/handler"
	"github.com/NikolosHGW/goph-keeper/internal/server/infrastructure/config"
	"github.com/NikolosHGW/goph-keeper/internal/server/infrastructure/db"
	"github.com/NikolosHGW/goph-keeper/internal/server/infrastructure/repository"
	"github.com/NikolosHGW/goph-keeper/internal/server/service"
	"github.com/NikolosHGW/goph-keeper/internal/server/usecase"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(fmt.Errorf("не удалось запустить сервер: %w", err))
	}
}

func run() error {
	config := config.NewConfig()

	myLogger, err := logger.NewLogger("info")
	if err != nil {
		return fmt.Errorf("не удалось инициализировать логгер: %w", err)
	}

	database, err := db.InitDB(config.GetDatabaseURI(), &db.DBConnector{}, &db.Migrator{})
	if err != nil {
		return fmt.Errorf("не удалось инициализировать базу данных: %w", err)
	}

	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			myLogger.LogInfo("ошибка при закрытии базы данных: ", closeErr)
		}
	}()

	userRepo := repository.NewUser(database, myLogger)

	registerService := service.NewRegister(myLogger)
	tokenService := service.NewToken(myLogger, config.GetSecretKey())

	registerUsecase := usecase.NewRegister(registerService, tokenService, userRepo)

	listen, err := net.Listen("tcp", config.GetRunAddress())
	if err != nil {
		return fmt.Errorf("не удалось прослушать TCP: %w", err)
	}

	srv := grpc.NewServer()

	registerpb.RegisterRegisterServer(srv, handler.NewRegisterServer(registerUsecase))

	errChan := make(chan error, 1)

	go func() {
		myLogger.LogStringInfo("Запуск сервера", "address", config.GetRunAddress())
		if err := srv.Serve(listen); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errChan <- fmt.Errorf("ошибка при запуске сервера: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		myLogger.LogStringInfo("Получен сигнал завершения, отключаем сервер...", "address", config.GetRunAddress())
	case err := <-errChan:
		myLogger.LogInfo("Сервер завершился с ошибкой: ", err)

		return fmt.Errorf("горутина с запуском сервера вернула ошибку: %w", err)
	}

	srv.GracefulStop()

	myLogger.LogStringInfo("Сервер успешно остановлен", "address", config.GetRunAddress())

	return nil
}
