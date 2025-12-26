package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/FooxyS/todo-service/internal/domain"
)

func TestCreateTodo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		savedID := 1
		mockRepo := &TodoRepositoryMock{
			SaveFunc: func(ctx context.Context, todo domain.Todo) (int, error) {
				return savedID, nil
			},
		}

		ctx := context.Background()

		inputTitle := "complete the game"
		inputTodo := domain.Todo{
			Title: inputTitle,
		}

		usecase := New(mockRepo)

		// act
		id, err := usecase.CreateTodo(ctx, inputTodo)

		// assert
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		if id != savedID {
			t.Errorf("unexpectedID: got %d, want %d", id, savedID)
		}

		wantCalls := 1
		if mockRepo.SaveCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.SaveCalls, wantCalls)
		}

		if mockRepo.LastSavedTodo.Title != inputTodo.Title {
			t.Errorf("unexpected title: got %q, want %q", mockRepo.LastSavedTodo.Title, inputTodo.Title)
		}
	})

	t.Run("failed to validate -> error", func(t *testing.T) {
		// preparing
		mockRepo := &TodoRepositoryMock{}

		usecase := TodoUseCase{
			TodoRepo: mockRepo,
		}

		ctx := context.Background()

		inputTodo := domain.Todo{}

		// act
		_, err := usecase.CreateTodo(ctx, inputTodo)

		// assert
		if !errors.Is(err, domain.ErrNoTitle) {
			t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrNoTitle)
		}

		wantCalls := 0
		if mockRepo.SaveCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.SaveCalls, wantCalls)
		}
	})

	t.Run("failed to save todo -> error", func(t *testing.T) {
		// preparing
		expectedError := errors.New("failed to save todo")

		mockRepo := &TodoRepositoryMock{
			SaveFunc: func(ctx context.Context, todo domain.Todo) (int, error) {
				return 0, expectedError
			},
		}

		usecase := TodoUseCase{
			TodoRepo: mockRepo,
		}

		ctx := context.Background()

		inputTodo := domain.Todo{
			Title: "read the book",
		}

		// act
		_, err := usecase.CreateTodo(ctx, inputTodo)
		if !errors.Is(err, expectedError) {
			t.Fatalf("unexpected error: got %v, want %v", err, expectedError)
		}

		wantCalls := 1
		if mockRepo.SaveCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.SaveCalls, wantCalls)
		}
	})
}

func TestGetAllTodos(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		// todos data
		savedFirstID := 1
		savedFirstTitle := "read the book"

		savedSecondID := 2
		savedSecondTitle := "complete the game"

		inputTodos := []domain.Todo{
			{ID: savedFirstID, Title: savedFirstTitle},
			{ID: savedSecondID, Title: savedSecondTitle},
		}

		mockRepo := &TodoRepositoryMock{
			ReadAllFunc: func(ctx context.Context) ([]domain.Todo, error) {
				return inputTodos, nil
			},
		}

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		todos, err := usecase.GetAllTodos(ctx)

		// assert
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		if len(todos) != len(inputTodos) {
			t.Fatalf("unexpected length: got %d, want %d", len(todos), len(inputTodos))
		}

		for i, todo := range todos {
			if todo.ID != inputTodos[i].ID {
				t.Errorf("unexpected id: got %d, want %d", todo.ID, inputTodos[i].ID)
			}

			if todo.Title != inputTodos[i].Title {
				t.Errorf("unexpected title: got %q, want %q", todo.Title, inputTodos[i].Title)
			}
		}

		wantCalls := 1
		if mockRepo.ReadAllCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.ReadAllCalls, wantCalls)
		}
	})

	t.Run("DB faield -> error", func(t *testing.T) {
		// preparing
		returnedError := errors.New("some error in DB")

		mockRepo := &TodoRepositoryMock{
			ReadAllFunc: func(ctx context.Context) ([]domain.Todo, error) {
				return []domain.Todo{}, returnedError
			},
		}

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		_, err := usecase.GetAllTodos(ctx)

		// assert
		if !errors.Is(err, returnedError) {
			t.Fatalf("unexpected error: got %v, want %v", err, returnedError)
		}

		wantCalls := 1
		if mockRepo.ReadAllCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.ReadAllCalls, wantCalls)
		}
	})
}

func TestGetTodoByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		savedID := 1
		savedTitle := "pre order the resident evil requiem"
		savedTodo := domain.Todo{
			ID:    savedID,
			Title: savedTitle,
		}

		mockRepo := &TodoRepositoryMock{
			GetByIDFunc: func(ctx context.Context, id int) (domain.Todo, error) {
				return savedTodo, nil
			},
		}

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		todo, err := usecase.GetTodoByID(ctx, savedID)

		// assert
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		if todo.ID != savedID {
			t.Errorf("unexpected id: got %d, want %d", todo.ID, savedID)
		}

		if todo.Title != savedTitle {
			t.Errorf("unexpected title: got %q, want %q", todo.Title, savedTitle)
		}

		wantCalls := 1
		if mockRepo.GetByIDCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.GetByIDCalls, wantCalls)
		}

		if mockRepo.LastGetID != savedID {
			t.Errorf("unexpected id arg: got %d, want %d", mockRepo.LastGetID, savedID)
		}
	})

	t.Run("todo does not exist -> error", func(t *testing.T) {
		// preparing
		inputId := 1

		mockRepo := &TodoRepositoryMock{
			GetByIDFunc: func(ctx context.Context, id int) (domain.Todo, error) {
				return domain.Todo{}, domain.ErrTodoNotExist
			},
		}

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		_, err := usecase.GetTodoByID(ctx, inputId)

		// assert
		if !errors.Is(err, domain.ErrTodoNotExist) {
			t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrTodoNotExist)
		}

		wantCalls := 1
		if mockRepo.GetByIDCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.GetByIDCalls, wantCalls)
		}
	})
}

func TestUpdateTodoByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		mockRepo := &TodoRepositoryMock{
			UpdateByIDFunc: func(ctx context.Context, id int, todo domain.Todo) error {
				return nil
			},
		}

		inputTodo := domain.Todo{
			Title: "complete the game",
		}

		inputID := 1

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		err := usecase.UpdateTodoByID(ctx, inputID, inputTodo)

		// assert
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		if mockRepo.LastSavedTodo.Title != inputTodo.Title {
			t.Errorf("unexpected todo arg: got %q , want %q", mockRepo.LastSavedTodo.Title, inputTodo.Title)
		}

		if mockRepo.LastGetID != inputID {
			t.Errorf("unexpected id arg: got %d , want %d", mockRepo.LastGetID, inputID)
		}

		wantCalls := 1
		if mockRepo.UpdateByIDCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.UpdateByIDCalls, wantCalls)
		}
	})

	t.Run("failed to validate -> error", func(t *testing.T) {
		// preparing
		inputID := 1

		inputTodo := domain.Todo{}

		mockRepo := &TodoRepositoryMock{}

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		err := usecase.UpdateTodoByID(ctx, inputID, inputTodo)

		// assert
		if !errors.Is(err, domain.ErrNoTitle) {
			t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrNoTitle)
		}

		wantCalls := 0
		if mockRepo.UpdateByIDCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.UpdateByIDCalls, wantCalls)
		}
	})

	t.Run("failed to update -> error", func(t *testing.T) {
		// preparing
		mockRepo := &TodoRepositoryMock{
			UpdateByIDFunc: func(ctx context.Context, id int, todo domain.Todo) error {
				return domain.ErrTodoNotExist
			},
		}

		inputTodo := domain.Todo{
			Title: "complete the game",
		}

		inputID := 1

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		err := usecase.UpdateTodoByID(ctx, inputID, inputTodo)

		// assert
		if !errors.Is(err, domain.ErrTodoNotExist) {
			t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrTodoNotExist)
		}

		if mockRepo.LastSavedTodo.Title != inputTodo.Title {
			t.Errorf("unexpected todo arg: got %q , want %q", mockRepo.LastSavedTodo.Title, inputTodo.Title)
		}

		if mockRepo.LastGetID != inputID {
			t.Errorf("unexpected id arg: got %d , want %d", mockRepo.LastGetID, inputID)
		}

		wantCalls := 1
		if mockRepo.UpdateByIDCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.UpdateByIDCalls, wantCalls)
		}
	})
}

func TestDeleteTodoByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		mockRepo := &TodoRepositoryMock{
			DeleteByIDFunc: func(ctx context.Context, id int) error {
				return nil
			},
		}

		inputID := 1

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		err := usecase.DeleteTodoByID(ctx, inputID)

		// assert
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		if mockRepo.LastGetID != inputID {
			t.Errorf("unexpecnted id arg: got %d, want %d", mockRepo.LastGetID, inputID)
		}

		wantCalls := 1
		if mockRepo.DeleteByIDCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.DeleteByIDCalls, wantCalls)
		}
	})

	t.Run("DB failed -> error", func(t *testing.T) {
		// preparing
		returnedError := errors.New("some error in DB")

		mockRepo := &TodoRepositoryMock{
			DeleteByIDFunc: func(ctx context.Context, id int) error {
				return returnedError
			},
		}

		inputID := 1

		ctx := context.Background()

		usecase := New(mockRepo)

		// act
		err := usecase.DeleteTodoByID(ctx, inputID)

		// assert
		if !errors.Is(err, returnedError) {
			t.Fatalf("unexpected error: got %v, want %v", err, returnedError)
		}

		if mockRepo.LastGetID != inputID {
			t.Errorf("unexpected id arg: got %d, want %d", mockRepo.LastGetID, inputID)
		}

		wantCalls := 1
		if mockRepo.DeleteByIDCalls != wantCalls {
			t.Errorf("unexpected calls: got %d, want %d", mockRepo.DeleteByIDCalls, wantCalls)
		}
	})
}
