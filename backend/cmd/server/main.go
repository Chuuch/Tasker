package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chuuch/Tasker/backend/internal/config"
	"github.com/Chuuch/Tasker/backend/internal/database"
	"github.com/Chuuch/Tasker/backend/internal/http/router"
	"github.com/Chuuch/Tasker/backend/internal/logger"
	"github.com/Chuuch/Tasker/backend/internal/metrics"
	"github.com/Chuuch/Tasker/backend/internal/task/handler"
	"github.com/Chuuch/Tasker/backend/internal/task/repository"
	"github.com/Chuuch/Tasker/backend/internal/task/service"
)

func main() {

	log := logger.New()
	cfg := config.Load()
	metrics.Register()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := database.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	taskRepository := repository.NewPostgresTaskRepository(db)
	taskService := service.NewTaskService(taskRepository)
	taskHandler := handler.NewTaskHandler(taskService)

	router := router.NewRouter(taskHandler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info("API listnening", "address", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		log.Error("server error", "error", err)

	case sig := <-shutdown:
		log.Info(
			"shutdown signal received",
			"signal", sig.String(),
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("gracefull shutdown failed", "error", err)
		}
	}

	log.Info("server stopped")
}
