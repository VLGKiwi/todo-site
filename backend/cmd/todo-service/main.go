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

	"github.com/VLGKiwi/todo-site/backend/internal/adapter/memory"  // Твоя БД
	"github.com/VLGKiwi/todo-site/backend/internal/controller/rest" // Твои контроллеры
	"github.com/VLGKiwi/todo-site/backend/internal/usecase"         // Твои use cases
)

func main() {
	// LOGGER
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// DB - in-memory (данные будут теряться при перезапуске!)
	db := memory.New()

	// USECASE
	uc := usecase.New(db)

	// SERVER
	router := rest.NewRouter(uc)

	// Добавляем CORS middleware
	corsRouter := addCorsMiddleware(router)

	// Добавляем health check endpoint
	mux := http.NewServeMux()
	mux.Handle("/", corsRouter)
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// Все API запросы через CORS middleware
		corsRouter.ServeHTTP(w, r)
	})

	// Получаем порт из переменных окружения (Render автоматически устанавливает PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Канал для ошибок сервера
	errCh := make(chan error, 1)

	slog.Info("Starting server", "addr", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errCh:
		slog.Error("Server failed to start", "error", err)
		return
	case sig := <-sigCh:
		slog.Info("Signal received", "signal", sig)
	}

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
	}

	slog.Info("Server stopped")
}

// Health check endpoint для Render
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// CORS middleware
func addCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любых источников (можно ограничить)
		origin := r.Header.Get("Origin")

		// Для разработки - разрешаем все
		// В продакшене лучше указать конкретные домены
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обрабатываем preflight запросы
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
