package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FooxyS/todo-service/internal/adapter/memory"
	"github.com/FooxyS/todo-service/internal/controller/rest"
	"github.com/FooxyS/todo-service/internal/usecase"
)

func main() {
	// LOGGER
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// DB
	db := memory.New()

	// USECASE
	uc := usecase.New(db)

	// SERVER
	router := rest.NewRouter(uc)
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// create error channel to receive server errors
	errCh := make(chan error, 1)

	slog.Info("start server", "addr", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errCh:
		slog.Error("server failed to start", "error", err)
		return
	case sig := <-sigCh:
		slog.Info("signal received", "signal", sig)
	}

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("http server shutdown error", "error", err)
	}

	slog.Info("server stopped")
}
