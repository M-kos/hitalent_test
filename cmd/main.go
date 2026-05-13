package main

import (
	"log/slog"

	"github.com/M-kos/hitalent_test/internal/config"
	"github.com/M-kos/hitalent_test/pkg/logger"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err.Error())
		return
	}

	log := logger.New(conf)
}
