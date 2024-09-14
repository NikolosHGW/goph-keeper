package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NikolosHGW/goph-keeper/internal/handler"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/config"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/db"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/repository"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/router"
	"github.com/NikolosHGW/goph-keeper/internal/service"
	"github.com/NikolosHGW/goph-keeper/internal/usecase"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
)

const shutdownTime = 5

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

	userService := service.NewUser(myLogger)

	registerUsecase := usecase.NewRegisterUser(userService, userRepo)

	handlers := &handler.Handlers{
		RegisterHandler: handler.NewRegisterHandler(registerUsecase, myLogger, config.GetSecretKey()),
	}

	r := router.NewRouter(handlers)

	srv := &http.Server{
		Addr:    config.GetRunAddress(),
		Handler: r,
	}

	errChan := make(chan error, 1)

	go func() {
		myLogger.LogStringInfo("Запуск сервера", "address", config.GetRunAddress())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), shutdownTime*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("ошибка при завершении работы сервера: %w", err)
	}

	myLogger.LogStringInfo("Сервер успешно остановлен", "address", config.GetRunAddress())

	return nil
}
