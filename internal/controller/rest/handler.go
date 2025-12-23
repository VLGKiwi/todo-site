package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/FooxyS/todo-service/internal/domain"
)

type UseCase interface {
	CreateTodo(ctx context.Context, todo domain.Todo) (int, error)
	GetAllTodos(ctx context.Context) []domain.Todo
	GetTodoByID(ctx context.Context, id int) (domain.Todo, error)
	UpdateTodoByID(ctx context.Context, todo domain.Todo) error
	DeleteTodoByID(ctx context.Context, id int) error
}

type Handlers struct {
	UseCase UseCase
}

func NewHandlers(usecase UseCase) *Handlers {
	return &Handlers{
		UseCase: usecase,
	}
}

func (h *Handlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo domain.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		// TODO: add logging of error with slog
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := h.UseCase.CreateTodo(r.Context(), todo)
	if errors.Is(err, domain.ErrNoTitle) {
		// TODO: add log
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		// TODO: log
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// set response headers and status code
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/todos/%d", id))

	w.WriteHeader(http.StatusCreated)

	resp := map[string]int{"id": id}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// log
	}
}

// TODO: implement other handlers
