package main

import (
	"fmt"
	"log"

	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/config"
	"github.com/NikolosHGW/goph-keeper/internal/infrastructure/db"
	"github.com/NikolosHGW/goph-keeper/pkg/logger"
	"go.uber.org/zap"
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
			myLogger.Fatal("ошибка при закрытии базы данных: ", zap.Error(err))
		}
	}()

	return nil
}
