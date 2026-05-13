package main

import (
	"log/slog"

	"github.com/M-kos/hitalent_test/internal/config"
	"github.com/M-kos/hitalent_test/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err.Error())
		return
	}

	log := logger.New(conf)

	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
