package usecase

import (
	"context"

	"github.com/FooxyS/todo-service/internal/domain"
)

type TodoRepositoryMock struct {
	SaveFunc       func(ctx context.Context, todo domain.Todo) (int, error)
	GetByIDFunc    func(ctx context.Context, id int) (domain.Todo, error)
	UpdateByIDFunc func(ctx context.Context, id int, todo domain.Todo) error
	DeleteByIDFunc func(ctx context.Context, id int) error
	ReadAllFunc    func(ctx context.Context) ([]domain.Todo, error)

	SaveCalls       int
	GetByIDCalls    int
	UpdateByIDCalls int
	DeleteByIDCalls int
	ReadAllCalls    int

	LastSavedTodo domain.Todo
	LastGetID     int
}

func (t *TodoRepositoryMock) Save(ctx context.Context, todo domain.Todo) (int, error) {
	t.SaveCalls++
	t.LastSavedTodo = todo

	if t.SaveFunc == nil {
		panic("SaveFunc is nil")
	}

	return t.SaveFunc(ctx, todo)
}

func (t *TodoRepositoryMock) GetByID(ctx context.Context, id int) (domain.Todo, error) {
	t.GetByIDCalls++
	t.LastGetID = id

	if t.GetByIDFunc == nil {
		panic("GetByIDFunc is nil")
	}

	return t.GetByIDFunc(ctx, id)
}

func (t *TodoRepositoryMock) UpdateByID(ctx context.Context, id int, todo domain.Todo) error {
	t.UpdateByIDCalls++
	t.LastSavedTodo = todo
	t.LastGetID = id

	if t.UpdateByIDFunc == nil {
		panic("UpdateByIDFunc is nil")
	}

	return t.UpdateByIDFunc(ctx, id, todo)
}

func (t *TodoRepositoryMock) DeleteByID(ctx context.Context, id int) error {
	t.DeleteByIDCalls++
	t.LastGetID = id

	if t.DeleteByIDFunc == nil {
		panic("DeleteByIDFunc is nil")
	}

	return t.DeleteByIDFunc(ctx, id)
}

func (t *TodoRepositoryMock) ReadAll(ctx context.Context) ([]domain.Todo, error) {
	t.ReadAllCalls++

	if t.ReadAllFunc == nil {
		panic("ReadAllFunc is nil")
	}

	return t.ReadAllFunc(ctx)
}
