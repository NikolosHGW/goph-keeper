package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolosHGW/goph-keeper/internal/handler"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/config"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/db"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/repository"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/router"
	"github.com/NikolosHGW/goph-keeper/internal/service"
	"github.com/NikolosHGW/goph-keeper/internal/usecase"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
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

	database, err := db.InitDB(config.GetDatabaseURI())
	if err != nil {
		return fmt.Errorf("не удалось инициализировать базу данных: %w", err)
	}

	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			myLogger.LogInfo("ошибка при закрытии базы данных: ", err)
		}
	}()

	userRepo := repository.NewUser(database, myLogger)

	userService := service.NewUser(myLogger)

	registerUsecase := usecase.NewRegisterUser(userService, userRepo)

	handlers := &handler.Handlers{
		AuthHandler: handler.NewAuthHandler(registerUsecase, userService, myLogger, config.GetSecretKey()),
	}

	r := router.NewRouter(handlers)

	myLogger.LogStringInfo("Running server", "address", config.GetRunAddress())

	err = http.ListenAndServe(config.GetRunAddress(), r)

	if err != nil {
		return fmt.Errorf("ошибка при запуске сервера: %w", err)
	}

	return nil
}
