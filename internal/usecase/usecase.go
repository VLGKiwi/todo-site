package Todousecase

import (
	"fmt"

	"github.com/FooxyS/todo-service/internal/domain"
)

type TodoRepository interface {
	Save(domain.Todo) (int, error)
	GetByID(id int) (domain.Todo, error)
	UpdateByID(domain.Todo) error
	DeleteByID(id int) error
	ReadAll() []domain.Todo
}

type TodoUseCase struct {
	todoRepo TodoRepository
}

func (u *TodoUseCase) CreateTodo(todo domain.Todo) (int, error) {
	// validate todo
	if err := todo.Validate(); err != nil {
		return 0, fmt.Errorf("validate todo: %w", err)
	}

	// save todo in db
	id, err := u.todoRepo.Save(todo)
	if err != nil {
		return 0, fmt.Errorf("save todo in db: %w", err)
	}

	return id, nil
}

func (u *TodoUseCase) GetAllTodos() []domain.Todo {
	// get all todos
	return u.todoRepo.ReadAll()
}

func (u *TodoUseCase) GetTodoByID(id int) (domain.Todo, error) {
	todo, err := u.todoRepo.GetByID(id)
	if err != nil {
		return domain.Todo{}, fmt.Errorf("get todo by id: %w", err)
	}
	return todo, nil
}

func (u *TodoUseCase) UpdateTodoByID(todo domain.Todo) error {
	// validate todo
	if err := todo.Validate(); err != nil {
		return fmt.Errorf("validate todo: %w", err)
	}

	// update todo in db
	if err := u.todoRepo.UpdateByID(todo); err != nil {
		return fmt.Errorf("update todo in db: %w", err)
	}

	return nil
}

func (u *TodoUseCase) DeleteTodoByID(id int) error {
	return u.todoRepo.DeleteByID(id)
}
