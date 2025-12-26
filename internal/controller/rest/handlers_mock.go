package rest

import (
	"context"

	"github.com/FooxyS/todo-service/internal/domain"
)

type UseCaseMock struct {
	CreateTodoFunc     func(ctx context.Context, todo domain.Todo) (int, error)
	GetAllTodosFunc    func(ctx context.Context) ([]domain.Todo, error)
	GetTodoByIDFunc    func(ctx context.Context, id int) (domain.Todo, error)
	UpdateTodoByIDFunc func(ctx context.Context, id int, todo domain.Todo) error
	DeleteTodoByIDFunc func(ctx context.Context, id int) error

	CreateTodoCalls     int
	GetAllTodosCalls    int
	GetTodoByIDCalls    int
	UpdateTodoByIDCalls int
	DeleteTodoByIDCalls int

	LastSavedTodo domain.Todo
	LastGetID     int
}

func (u *UseCaseMock) CreateTodo(ctx context.Context, todo domain.Todo) (int, error) {
	u.LastSavedTodo = todo
	u.CreateTodoCalls++

	if u.CreateTodoFunc == nil {
		panic("CreateTodoFunc is nil")
	}

	return u.CreateTodoFunc(ctx, todo)
}

func (u *UseCaseMock) GetAllTodos(ctx context.Context) ([]domain.Todo, error) {
	u.GetAllTodosCalls++
	if u.GetAllTodosFunc == nil {
		panic("GetAllTodosFunc is nil")
	}

	return u.GetAllTodosFunc(ctx)
}

func (u *UseCaseMock) GetTodoByID(ctx context.Context, id int) (domain.Todo, error) {
	u.LastGetID = id
	u.GetTodoByIDCalls++

	if u.GetTodoByIDFunc == nil {
		panic("GetTodoByIDFunc is nil")
	}

	return u.GetTodoByIDFunc(ctx, id)
}

func (u *UseCaseMock) UpdateTodoByID(ctx context.Context, id int, todo domain.Todo) error {
	u.LastGetID = id
	u.LastSavedTodo = todo
	u.UpdateTodoByIDCalls++

	if u.UpdateTodoByIDFunc == nil {
		panic("UpdateTodoByIDFunc is nil")
	}

	return u.UpdateTodoByIDFunc(ctx, id, todo)
}

func (u *UseCaseMock) DeleteTodoByID(ctx context.Context, id int) error {
	u.LastGetID = id
	u.DeleteTodoByIDCalls++

	if u.DeleteTodoByIDFunc == nil {
		panic("DeleteTodoByIDFunc is nil")
	}

	return u.DeleteTodoByIDFunc(ctx, id)
}
