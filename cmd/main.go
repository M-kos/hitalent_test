package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/M-kos/hitalent_test/internal/config"
	"github.com/M-kos/hitalent_test/internal/db"
	"github.com/M-kos/hitalent_test/internal/handler"
	"github.com/M-kos/hitalent_test/internal/repository"
	"github.com/M-kos/hitalent_test/internal/service"
	"github.com/M-kos/hitalent_test/pkg/logger"
)

const (
	shutdownTimeout = 15 * time.Second
)

func main() {
	conf, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err.Error())
		return
	}

	log := logger.New(conf)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	gormDb, err := db.OpenGormConnection(ctx, conf)
	if err != nil {
		log.Error("failed to connect to database", "error", err.Error())
		return
	}

	repo := repository.NewDepartmentRepository(gormDb)
	depService := service.NewDepartmentService(repo)

	router := http.NewServeMux()
	handler.NewDepartmentHandler(router, conf, log, depService)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("listen and serve", "msg", err.Error())
			cancel()
		}
	}()

	log.Info("server starting", slog.Int("port", conf.Port))

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown server error", "msg", err.Error())
	}

	sqlDB, err := gormDb.DB()
	if err != nil {
		log.Error("failed to get sql db from gorm", "error", err.Error())
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Error("failed to close db", "error", err.Error())
		}
	}

	log.Info("server shutdown", slog.Int("port", conf.Port))
}
