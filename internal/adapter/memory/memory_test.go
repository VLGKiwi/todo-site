package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/FooxyS/todo-service/internal/domain"
)

func TestSave(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx := context.Background()

		inputTitle := "read the book"
		todo := domain.Todo{
			Title: inputTitle,
		}

		// act
		id, err := todoRepo.Save(ctx, todo)

		// assert
		if err != nil {
			t.Fatalf("got %v, want nil", err)
		}

		// check saved todo in db
		wantId := 1
		if id != wantId {
			t.Fatalf("must return correct id: got %d, want %d", id, wantId)
		}

		savedTodo := todoRepo.DB[id]
		if err != nil {
			t.Fatalf("todo not saved: got %v, want nil", err)
		}

		if savedTodo.ID != wantId {
			t.Fatalf("must put id into saving todo: got %d, want %d", id, wantId)
		}

		wantTitle := inputTitle
		if savedTodo.Title != wantTitle {
			t.Fatalf("incorrect title: got %s, want %s", savedTodo.Title, wantTitle)
		}

		wantNextID := 2
		if todoRepo.NextID != wantNextID {
			t.Fatalf("must increment nextID field: got %d, want %d", todoRepo.NextID, wantNextID)
		}
	})

	t.Run("context canceled -> throw error", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		todo := domain.Todo{
			Title: "read the book",
		}

		// act
		id, err := todoRepo.Save(ctx, todo)

		// assert
		wantErr := ctx.Err()
		if !errors.Is(err, wantErr) {
			t.Fatalf("gotErr %v, want %v", err, wantErr)
		}

		// check that id was not incremented
		wantNExtID := 1
		if todoRepo.NextID != wantNExtID {
			t.Errorf("must not increment when error: got %d, want %d", todoRepo.NextID, wantNExtID)
		}

		// check that todo was not saved to DB
		_, ok := todoRepo.DB[id]
		if ok {
			t.Errorf("must not save todo in DB when error")
		}
	})

}

func TestGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx := context.Background()

		inputID := 1
		inputTitle := "read the book"
		inputTodo := domain.Todo{
			ID:    inputID,
			Title: inputTitle,
		}

		todoRepo.DB[inputID] = inputTodo

		// act
		todo, err := todoRepo.GetByID(ctx, inputID)

		// assert
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		if todo.ID != inputID {
			t.Errorf("unexpected todoID: got %d, want %d", todo.ID, inputID)
		}

		if todo.Title != inputTitle {
			t.Errorf("unexpected todo's title: got %s, want %s", todo.Title, inputTitle)
		}
	})

	t.Run("context canceled -> error", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		inputID := 1

		// act
		_, err := todoRepo.GetByID(ctx, inputID)

		// assert
		if !errors.Is(err, ctx.Err()) {
			t.Fatalf("unexpected error: got %v, want %v", err, ctx.Err())
		}
	})

	t.Run("todoNotExist -> error", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx := context.Background()

		inputID := 1

		// act
		_, err := todoRepo.GetByID(ctx, inputID)

		// assert
		if !errors.Is(err, domain.ErrTodoNotExist) {
			t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrTodoNotExist)
		}
	})
}

func TestUpdateByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx := context.Background()

		// saved todo
		savedTitle := "read the book"
		savedID := 1
		savedTodo := domain.Todo{
			ID:    savedID,
			Title: savedTitle,
		}

		todoRepo.DB[savedID] = savedTodo

		// updated todo
		updatedTitle := "complete the god of war"
		updatedDescription := "amazing game!"
		updatedCompleted := true

		todoForUpdate := domain.Todo{
			Title:       updatedTitle,
			Description: updatedDescription,
			Completed:   updatedCompleted,
		}

		// act
		err := todoRepo.UpdateByID(ctx, savedID, todoForUpdate)

		// assert
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		savedUpdatedTodo := todoRepo.DB[savedID]
		if savedUpdatedTodo.ID != savedID {
			t.Errorf("unexpected id: got %d, want %d", savedUpdatedTodo.ID, savedID)
		}

		if savedUpdatedTodo.Title != updatedTitle {
			t.Errorf("unexpected title: got %s, want %s", savedUpdatedTodo.Title, updatedTitle)
		}

		if savedUpdatedTodo.Description != updatedDescription {
			t.Errorf("unexpected description: got %s, want %s", savedUpdatedTodo.Description, updatedDescription)
		}

		if savedUpdatedTodo.Completed != updatedCompleted {
			t.Errorf("unexpected completed status: got %t, want %t", savedUpdatedTodo.Completed, updatedCompleted)
		}
	})

	t.Run("context canceled -> error", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		inputID := 1

		inputTodo := domain.Todo{}

		// act
		err := todoRepo.UpdateByID(ctx, inputID, inputTodo)

		// assert
		if !errors.Is(err, ctx.Err()) {
			t.Fatalf("unexpected error: got %v, want %v", err, ctx.Err())
		}
	})

	t.Run("todoNotExist -> error", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx := context.Background()

		inputID := 1

		inputTodo := domain.Todo{}

		// act
		err := todoRepo.UpdateByID(ctx, inputID, inputTodo)

		// assert
		if !errors.Is(err, domain.ErrTodoNotExist) {
			t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrTodoNotExist)
		}
	})
}

func TestDeleteByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx := context.Background()

		// saved todo
		savedTitle := "read the book"
		savedID := 1
		savedTodo := domain.Todo{
			ID:    savedID,
			Title: savedTitle,
		}

		todoRepo.DB[savedID] = savedTodo

		// act
		err := todoRepo.DeleteByID(ctx, savedID)
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		_, ok := todoRepo.DB[savedID]
		if ok {
			t.Errorf("todo was not deleted")
		}
	})

	t.Run("context canceled -> error", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		inputID := 1

		// act
		err := todoRepo.DeleteByID(ctx, inputID)

		// assert
		if !errors.Is(err, ctx.Err()) {
			t.Fatalf("unexpected error: got %v, want %v", err, ctx.Err())
		}
	})

	t.Run("todoNotExist -> error", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx := context.Background()

		inputID := 1

		// act
		err := todoRepo.DeleteByID(ctx, inputID)

		// assert
		if !errors.Is(err, domain.ErrTodoNotExist) {
			t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrTodoNotExist)
		}
	})
}

func TestReadAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx := context.Background()

		// todos data
		savedFirstID := 1
		savedFirstTitle := "read the book"

		savedSecondID := 2
		savedSecondTitle := "complete the game"

		savedThirdID := 3
		savedThirdTitle := "get an internship at ecom.tech"

		savedTodos := []domain.Todo{
			{ID: savedFirstID, Title: savedFirstTitle},
			{ID: savedSecondID, Title: savedSecondTitle},
			{ID: savedThirdID, Title: savedThirdTitle},
		}

		// saving input todos
		for _, todo := range savedTodos {
			todoRepo.DB[todo.ID] = todo
		}

		// act
		getTodos, err := todoRepo.ReadAll(ctx)

		// assert
		if err != nil {
			t.Fatalf("unexpected error: got %v, want nil", err)
		}

		for i, todo := range getTodos {
			if todo.ID != savedTodos[i].ID {
				t.Errorf("unexpected id: got %d, want %d", todo.ID, savedTodos[i].ID)
			}

			if todo.Title != savedTodos[i].Title {
				t.Errorf("unexpected title: got %s, want %s", todo.Title, savedTodos[i].Title)
			}
		}
	})

	t.Run("context canceled -> error", func(t *testing.T) {
		// preparing
		todoRepo := New()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		_, err := todoRepo.ReadAll(ctx)

		// assert
		if !errors.Is(err, ctx.Err()) {
			t.Fatalf("unexpected error: got %v, want %v", err, ctx.Err())
		}
	})
}
