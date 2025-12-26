package rest

import "net/http"

func NewRouter(usecase UseCase) http.Handler {
	mux := http.NewServeMux()

	handlers := NewHandlers(usecase)

	mux.HandleFunc("POST /todos", handlers.CreateTodoHandler)
	mux.HandleFunc("GET /todos", handlers.GetAllTodosHandler)
	mux.HandleFunc("GET /todos/{id}", handlers.GetTodoHandler)
	mux.HandleFunc("PUT /todos/{id}", handlers.UpdateTodoHandler)
	mux.HandleFunc("DELETE /todos/{id}", handlers.DeleteTodoHandler)

	wrappedMux := LoggingMiddleware(mux)

	return wrappedMux
}
