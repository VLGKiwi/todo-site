package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FooxyS/todo-service/internal/domain"
)

type UseCase interface {
	CreateTodo(ctx context.Context, todo domain.Todo) (int, error)
	GetAllTodos(ctx context.Context) ([]domain.Todo, error)
	GetTodoByID(ctx context.Context, id int) (domain.Todo, error)
	UpdateTodoByID(ctx context.Context, id int, todo domain.Todo) error
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

func (h *Handlers) CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo domain.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		slog.Warn("failed to decode request", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := h.UseCase.CreateTodo(r.Context(), todo)
	if errors.Is(err, domain.ErrNoTitle) {
		slog.Warn("todo validation failed", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	} else if err != nil {
		slog.Error("failed to create todo", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// set response headers and status code
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/todos/%d", id))

	w.WriteHeader(http.StatusCreated)

	resp := map[string]int{"id": id}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to write response", "error", err)
	}

	slog.Info("todo created", "id", id)
}

func (h *Handlers) GetAllTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := h.UseCase.GetAllTodos(r.Context())
	if err != nil {
		slog.Error("failed to get all todos", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func (h *Handlers) GetTodoHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("failed to convert id from string", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	todo, err := h.UseCase.GetTodoByID(r.Context(), id)
	if errors.Is(err, domain.ErrTodoNotExist) {
		slog.Warn("failed to get todo by id", "error", err, "id", id)
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("failed to get todo by id", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(todo); err != nil {
		slog.Error("failed to encode response", "error", err)
	}

}

func (h *Handlers) UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("failed to convert id from string", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var todo domain.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		slog.Warn("failed to decode request", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := h.UseCase.UpdateTodoByID(r.Context(), id, todo); errors.Is(err, domain.ErrTodoNotExist) {
		slog.Warn("failed to update todo", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if errors.Is(err, domain.ErrNoTitle) {
		slog.Warn("todo validation failed", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	} else if err != nil {
		slog.Error("failed to update todo", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"message": "todo successfully updated"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode response", "error", err)
	}

}

func (h *Handlers) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("failed to convert id from string", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := h.UseCase.DeleteTodoByID(r.Context(), id); errors.Is(err, domain.ErrTodoNotExist) {
		slog.Warn("failed to delete todo", "error", err, "id", id)
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("failed to delete todo", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
