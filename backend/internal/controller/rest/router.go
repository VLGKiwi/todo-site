package rest

import "net/http"

func NewRouter(usecase UseCase) http.Handler {
	mux := http.NewServeMux()

	handlers := NewHandlers(usecase)

	mux.HandleFunc("POST /api/todos", handlers.CreateTodoHandler)
	mux.HandleFunc("GET /api/todos", handlers.GetAllTodosHandler)
	mux.HandleFunc("GET /api/todos/{id}", handlers.GetTodoHandler)
	mux.HandleFunc("PUT /api/todos/{id}", handlers.UpdateTodoHandler)
	mux.HandleFunc("DELETE /api/todos/{id}", handlers.DeleteTodoHandler)

	wrappedMux := LoggingMiddleware(mux)

	return wrappedMux
}
