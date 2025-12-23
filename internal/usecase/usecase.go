package usecase

import (
	"fmt"

	"context"

	"github.com/FooxyS/todo-service/internal/domain"
)

type TodoRepository interface {
	Save(ctx context.Context, todo domain.Todo) (int, error)
	GetByID(ctx context.Context, id int) (domain.Todo, error)
	UpdateByID(ctx context.Context, todo domain.Todo) error
	DeleteByID(ctx context.Context, id int) error
	ReadAll(ctx context.Context) ([]domain.Todo, error)
}

type TodoUseCase struct {
	TodoRepo TodoRepository
}

func New(repo TodoRepository) *TodoUseCase {
	return &TodoUseCase{
		TodoRepo: repo,
	}
}

func (u *TodoUseCase) CreateTodo(ctx context.Context, todo domain.Todo) (int, error) {
	// validate todo
	if err := todo.Validate(); err != nil {
		return 0, fmt.Errorf("validate todo: %w", err)
	}

	// save todo in db
	id, err := u.TodoRepo.Save(ctx, todo)
	if err != nil {
		return 0, fmt.Errorf("save todo in db: %w", err)
	}

	return id, nil
}

func (u *TodoUseCase) GetAllTodos(ctx context.Context) ([]domain.Todo, error) {
	// get all todos
	return u.TodoRepo.ReadAll(ctx)
}

func (u *TodoUseCase) GetTodoByID(ctx context.Context, id int) (domain.Todo, error) {
	todo, err := u.TodoRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Todo{}, fmt.Errorf("get todo by id: %w", err)
	}
	return todo, nil
}

func (u *TodoUseCase) UpdateTodoByID(ctx context.Context, todo domain.Todo) error {
	// validate todo
	if err := todo.Validate(); err != nil {
		return fmt.Errorf("validate todo: %w", err)
	}

	// update todo in db
	if err := u.TodoRepo.UpdateByID(ctx, todo); err != nil {
		return fmt.Errorf("update todo in db: %w", err)
	}

	return nil
}

func (u *TodoUseCase) DeleteTodoByID(ctx context.Context, id int) error {
	return u.TodoRepo.DeleteByID(ctx, id)
}
