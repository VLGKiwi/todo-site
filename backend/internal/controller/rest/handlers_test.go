package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VLGKiwi/todo-site/backend/internal/domain"
)

func TestCreateTodoHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   string

		usecaseFunc func(ctx context.Context, todo domain.Todo) (int, error)

		wantCode        int
		wantContentType string
		wantLocation    string
		wantBody        string

		wantID    int
		wantTitle string

		wantCalls int
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body:   `{"title": "complete the task"}`,
			usecaseFunc: func(ctx context.Context, todo domain.Todo) (int, error) {
				return 1, nil
			},
			wantCode:        http.StatusCreated,
			wantContentType: "application/json",
			wantLocation:    "/todos/1",
			wantID:          1,
			wantTitle:       "complete the task",
			wantCalls:       1,
		},
		{
			name:            "failed to decode body -> error",
			method:          http.MethodPost,
			body:            `"title": "complete the task"`,
			wantCode:        http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "bad request\n",
			wantCalls:       0,
		},
		{
			name:   "failed todo validation -> error",
			method: http.MethodPost,
			body:   "{}",
			usecaseFunc: func(ctx context.Context, todo domain.Todo) (int, error) {
				return 0, domain.ErrNoTitle
			},
			wantCode:        http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "bad request\n",
			wantTitle:       "",
			wantCalls:       1,
		},
		{
			name:   "failed to create todo -> error",
			method: http.MethodPost,
			body:   `{"title": "complete the task"}`,
			usecaseFunc: func(ctx context.Context, todo domain.Todo) (int, error) {
				return 0, errors.New("some error in usecase")
			},
			wantCode:        http.StatusInternalServerError,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "internal server error\n",
			wantTitle:       "complete the task",
			wantCalls:       1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// preparing
			reader := strings.NewReader(tc.body)

			req := httptest.NewRequest(tc.method, "/todos", reader)
			rec := httptest.NewRecorder()

			useCaseMock := &UseCaseMock{}

			if tc.wantCalls == 0 {
				useCaseMock.CreateTodoFunc = func(ctx context.Context, todo domain.Todo) (int, error) {
					t.Fatalf("CreateTodo must not be called")
					return 0, nil
				}
			} else {
				useCaseMock.CreateTodoFunc = tc.usecaseFunc
			}

			handlers := Handlers{
				UseCase: useCaseMock,
			}

			// act
			handlers.CreateTodoHandler(rec, req)

			// assert
			if rec.Code != tc.wantCode {
				t.Errorf("unexpected status code: got %d, want %d", rec.Code, tc.wantCode)
			}

			if tc.wantContentType != "" {
				gotCT := rec.Header().Get("Content-Type")
				if gotCT != tc.wantContentType {
					t.Errorf("unexpected Content-Type: got %q, want %q", gotCT, tc.wantContentType)
				}
			}

			if tc.wantLocation != "" {
				gotLoc := rec.Header().Get("Location")
				if gotLoc != tc.wantLocation {
					t.Errorf("unexpected location: got %q, want %q", gotLoc, tc.wantLocation)
				}
			}

			if tc.wantBody != "" {
				if rec.Body.String() != tc.wantBody {
					t.Errorf("unexpected body: got %q, want %q", rec.Body.String(), tc.wantBody)
				}
			}

			if tc.wantID != 0 {
				var resp struct {
					ID int `json:"id"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode json: %v, body=%q", err, rec.Body.String())
				}
				if resp.ID != tc.wantID {
					t.Errorf("unexpected id: got %d, want %d", resp.ID, tc.wantID)
				}
			}

			if useCaseMock.CreateTodoCalls != tc.wantCalls {
				t.Errorf("unexpected calls: got %d, want %d", useCaseMock.CreateTodoCalls, tc.wantCalls)
			}

			if tc.wantCalls > 0 {
				if useCaseMock.LastSavedTodo.Title != tc.wantTitle {
					t.Errorf("unexpected title: got %q, want %q", useCaseMock.LastSavedTodo.Title, tc.wantTitle)
				}
			}
		})
	}
}

func TestGetAllTodosHandler(t *testing.T) {
	todos := []domain.Todo{
		{ID: 1, Title: "complete the game"},
		{ID: 2, Title: "read the book"},
	}

	tests := []struct {
		name   string
		method string

		usecaseFunc func(ctx context.Context) ([]domain.Todo, error)

		wantCode        int
		wantContentType string
		wantBody        string

		wantTodos bool

		wantCalls int
	}{
		{
			name:   "success",
			method: http.MethodGet,
			usecaseFunc: func(ctx context.Context) ([]domain.Todo, error) {
				return todos, nil
			},
			wantCode:        http.StatusOK,
			wantContentType: "application/json",
			wantTodos:       true,
			wantCalls:       1,
		},
		{
			name:   "failed to get todos -> error",
			method: http.MethodGet,
			usecaseFunc: func(ctx context.Context) ([]domain.Todo, error) {
				return []domain.Todo{}, errors.New("some error in usecase")
			},
			wantCode:        http.StatusInternalServerError,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "internal server error\n",
			wantCalls:       1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// preparing
			req := httptest.NewRequest(tc.method, "/todos", nil)
			rec := httptest.NewRecorder()

			useCaseMock := &UseCaseMock{}

			if tc.wantCalls == 0 {
				useCaseMock.GetAllTodosFunc = func(ctx context.Context) ([]domain.Todo, error) {
					t.Fatalf("GetAllTodos must not be called")
					return []domain.Todo{}, nil
				}
			} else {
				useCaseMock.GetAllTodosFunc = tc.usecaseFunc
			}

			handlers := Handlers{
				UseCase: useCaseMock,
			}

			// act
			handlers.GetAllTodosHandler(rec, req)

			// assert
			if tc.wantCode != rec.Code {
				t.Errorf("unexpected status code: got %d, want %d", rec.Code, tc.wantCode)
			}

			if tc.wantContentType != "" {
				gotCT := rec.Header().Get("Content-Type")
				if gotCT != tc.wantContentType {
					t.Errorf("unexpected Content-Type: got %q, want %q", gotCT, tc.wantContentType)
				}
			}

			if tc.wantBody != "" {
				gotBody := rec.Body.String()
				if gotBody != tc.wantBody {
					t.Errorf("unexpected body: got %q, want %q", gotBody, tc.wantBody)
				}
			}

			if tc.wantTodos {
				var resp []domain.Todo

				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode json: %v, body=%q", err, rec.Body.String())
				}

				if len(resp) != len(todos) {
					t.Errorf("unexpected length: got %d, want %d", len(resp), len(todos))
				}

				for i, todo := range resp {
					if todo.ID != todos[i].ID {
						t.Errorf("unexpected id: got %d, want %d", todo.ID, todos[i].ID)
					}
					if todo.Title != todos[i].Title {
						t.Errorf("unexpected title: got %q, want %q", todo.Title, todos[i].Title)
					}
				}
			}

			if useCaseMock.GetAllTodosCalls != tc.wantCalls {
				t.Errorf("unexpected calls: got %d, want %d", useCaseMock.GetAllTodosCalls, tc.wantCalls)
			}
		})
	}
}

func TestGetTodoHandler(t *testing.T) {
	todo := domain.Todo{
		ID:    1,
		Title: "read the book",
	}

	tests := []struct {
		name   string
		method string

		usecaseFunc func(ctx context.Context, id int) (domain.Todo, error)

		wantCode        int
		wantContentType string
		wantBody        string

		pathValue string
		wantTodo  bool
		wantID    int

		wantCalls int
	}{
		{
			name:   "success",
			method: http.MethodGet,
			usecaseFunc: func(ctx context.Context, id int) (domain.Todo, error) {
				return todo, nil
			},
			wantCode:        http.StatusOK,
			wantContentType: "application/json",
			pathValue:       "1",
			wantTodo:        true,
			wantID:          1,
			wantCalls:       1,
		},
		{
			name:            "failed to get path value -> error",
			method:          http.MethodGet,
			wantCode:        http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "bad request\n",
			pathValue:       "abc",
			wantCalls:       0,
		},
		{
			name:   "todo not exits -> error",
			method: http.MethodGet,
			usecaseFunc: func(ctx context.Context, id int) (domain.Todo, error) {
				return domain.Todo{}, domain.ErrTodoNotExist
			},
			wantCode:        http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "todo not found\n",
			pathValue:       "1",
			wantID:          1,
			wantCalls:       1,
		},
		{
			name:   "internal server error -> error",
			method: http.MethodGet,
			usecaseFunc: func(ctx context.Context, id int) (domain.Todo, error) {
				return domain.Todo{}, errors.New("some error from usecase")
			},
			wantCode:        http.StatusInternalServerError,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "internal server error\n",
			pathValue:       "1",
			wantID:          1,
			wantCalls:       1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// preparing
			req := httptest.NewRequest(tc.method, "/todos/1", nil)
			req.SetPathValue("id", tc.pathValue)
			rec := httptest.NewRecorder()

			useCaseMock := &UseCaseMock{
				GetTodoByIDFunc: tc.usecaseFunc,
			}

			if tc.wantCalls == 0 {
				useCaseMock.GetTodoByIDFunc = func(ctx context.Context, id int) (domain.Todo, error) {
					t.Fatalf("GetTodoByID must not be called")
					return domain.Todo{}, nil
				}
			} else {
				useCaseMock.GetTodoByIDFunc = tc.usecaseFunc
			}

			handlers := Handlers{
				UseCase: useCaseMock,
			}

			// act
			handlers.GetTodoHandler(rec, req)

			// assert
			if tc.wantCode != rec.Code {
				t.Errorf("unexpected status code: got %d, want %d", rec.Code, tc.wantCode)
			}

			if tc.wantContentType != "" {
				gotCT := rec.Header().Get("Content-Type")
				if gotCT != tc.wantContentType {
					t.Errorf("unexpected Content-Type: got %q, want %q", gotCT, tc.wantContentType)
				}
			}

			if tc.wantBody != "" {
				gotBody := rec.Body.String()
				if gotBody != tc.wantBody {
					t.Errorf("unexpected body: got %q, want %q", gotBody, tc.wantBody)
				}
			}

			if tc.wantTodo {
				var resp domain.Todo

				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode json: %v, body=%q", err, rec.Body.String())
				}

				if resp.ID != todo.ID {
					t.Errorf("unexpected id: got %d, want %d", resp.ID, todo.ID)
				}
				if resp.Title != todo.Title {
					t.Errorf("unexpected title: got %q, want %q", resp.Title, todo.Title)
				}
			}

			if useCaseMock.GetTodoByIDCalls != tc.wantCalls {
				t.Errorf("unexpected calls: got %d, want %d", useCaseMock.GetTodoByIDCalls, tc.wantCalls)
			}

			if tc.wantCalls > 0 {
				if useCaseMock.LastGetID != tc.wantID {
					t.Errorf("unexpected id: got %d, want %d", useCaseMock.LastGetID, tc.wantID)
				}
			}
		})
	}
}

func TestUpdateTodoHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   string

		usecaseFunc     func(ctx context.Context, id int, todo domain.Todo) error
		wantCode        int
		wantContentType string
		wantBody        string

		pathValue string
		wantTitle string
		wantID    int
		wantJson  string

		wantCalls int
	}{
		{
			name:   "success",
			method: http.MethodPut,
			body:   `{"title": "read the book"}`,
			usecaseFunc: func(ctx context.Context, id int, todo domain.Todo) error {
				return nil
			},
			wantCode:        http.StatusOK,
			wantContentType: "application/json",
			pathValue:       "1",
			wantTitle:       "read the book",
			wantID:          1,
			wantJson:        "todo successfully updated",
			wantCalls:       1,
		},
		{
			name:            "failed to get path value -> error",
			method:          http.MethodPut,
			wantCode:        http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "bad request\n",
			pathValue:       "abc",
			wantCalls:       0,
		},
		{
			name:            "failed to decode body -> error",
			method:          http.MethodPut,
			body:            `}`,
			wantCode:        http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "bad request\n",
			pathValue:       "1",
			wantCalls:       0,
		},
		{
			name:   "todo not exist -> error",
			method: http.MethodPut,
			body:   `{"title": "read the book"}`,
			usecaseFunc: func(ctx context.Context, id int, todo domain.Todo) error {
				return domain.ErrTodoNotExist
			},
			wantCode:        http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
			pathValue:       "1",
			wantTitle:       "read the book",
			wantID:          1,
			wantCalls:       1,
		},
		{
			name:   "failed to validate -> error",
			method: http.MethodPut,
			body:   `{}`,
			usecaseFunc: func(ctx context.Context, id int, todo domain.Todo) error {
				return domain.ErrNoTitle
			},
			wantCode:        http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
			pathValue:       "1",
			wantTitle:       "",
			wantID:          1,
			wantCalls:       1,
		},
		{
			name:   "internal server error -> error",
			method: http.MethodPut,
			body:   `{"title": "read the book"}`,
			usecaseFunc: func(ctx context.Context, id int, todo domain.Todo) error {
				return errors.New("some error from usecase")
			},
			wantCode:        http.StatusInternalServerError,
			wantContentType: "text/plain; charset=utf-8",
			pathValue:       "1",
			wantTitle:       "read the book",
			wantID:          1,
			wantCalls:       1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// preparing
			reader := strings.NewReader(tc.body)

			req := httptest.NewRequest(tc.method, "/todos/1", reader)
			req.SetPathValue("id", tc.pathValue)
			rec := httptest.NewRecorder()

			useCaseMock := &UseCaseMock{
				UpdateTodoByIDFunc: tc.usecaseFunc,
			}

			if tc.wantCalls == 0 {
				useCaseMock.UpdateTodoByIDFunc = func(ctx context.Context, id int, todo domain.Todo) error {
					t.Fatalf("UpdateTodoByID must not be called")
					return nil
				}
			} else {
				useCaseMock.UpdateTodoByIDFunc = tc.usecaseFunc
			}

			handlers := Handlers{
				UseCase: useCaseMock,
			}

			// act
			handlers.UpdateTodoHandler(rec, req)

			// assert
			if tc.wantCode != rec.Code {
				t.Errorf("unexpected status code: got %d, want %d", rec.Code, tc.wantCode)
			}

			if tc.wantContentType != "" {
				gotCT := rec.Header().Get("Content-Type")
				if gotCT != tc.wantContentType {
					t.Errorf("unexpected Content-Type: got %q, want %q", gotCT, tc.wantContentType)
				}
			}

			if tc.wantBody != "" {

				gotBody := rec.Body.String()
				if gotBody != tc.wantBody {
					t.Errorf("unexpected body: got %q, want %q", gotBody, tc.wantBody)
				}

			}

			if tc.wantJson != "" {
				var resp struct {
					Message string `json:"message"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode json: %v. body=%q", err, rec.Body.String())
				}
				if resp.Message != tc.wantJson {
					t.Errorf("unexpected json: got %q, want %q", resp.Message, tc.wantJson)
				}
			}

			if useCaseMock.UpdateTodoByIDCalls != tc.wantCalls {
				t.Errorf("unexpected calls: got %d, want %d", useCaseMock.UpdateTodoByIDCalls, tc.wantCalls)
			}

			if tc.wantCalls > 0 {
				if useCaseMock.LastGetID != tc.wantID {
					t.Errorf("unexpected id: got %d, want %d", useCaseMock.LastGetID, tc.wantID)
				}
				if useCaseMock.LastSavedTodo.Title != tc.wantTitle {
					t.Errorf("unexpected title: got %q, want %q", useCaseMock.LastSavedTodo.Title, tc.wantTitle)
				}
			}
		})
	}
}

func TestDeleteTodoHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string

		usecaseFunc     func(ctx context.Context, id int) error
		wantCode        int
		wantContentType string
		wantBody        string

		pathValue string
		wantID    int

		wantCalls int
	}{
		{
			name:   "success",
			method: http.MethodDelete,
			usecaseFunc: func(ctx context.Context, id int) error {
				return nil
			},
			wantCode:  http.StatusNoContent,
			wantBody:  "",
			pathValue: "1",
			wantID:    1,
			wantCalls: 1,
		},
		{
			name:            "failed to get path id",
			method:          http.MethodDelete,
			wantCode:        http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "bad request\n",
			pathValue:       "abc",
			wantCalls:       0,
		},
		{
			name:   "failed to delete -> error",
			method: http.MethodDelete,
			usecaseFunc: func(ctx context.Context, id int) error {
				return domain.ErrTodoNotExist
			},
			wantCode:        http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "todo not found\n",
			pathValue:       "1",
			wantID:          1,
			wantCalls:       1,
		},
		{
			name:   "internal server error -> error",
			method: http.MethodDelete,
			usecaseFunc: func(ctx context.Context, id int) error {
				return errors.New("some error from usecase")
			},
			wantCode:        http.StatusInternalServerError,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "internal server error\n",
			pathValue:       "1",
			wantID:          1,
			wantCalls:       1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// preparing
			req := httptest.NewRequest(tc.method, "/todos/1", nil)
			req.SetPathValue("id", tc.pathValue)
			rec := httptest.NewRecorder()

			useCaseMock := &UseCaseMock{
				DeleteTodoByIDFunc: tc.usecaseFunc,
			}

			if tc.wantCalls == 0 {
				useCaseMock.DeleteTodoByIDFunc = func(ctx context.Context, id int) error {
					t.Fatalf("DeleteTodoByID must not be called")
					return nil
				}
			} else {
				useCaseMock.DeleteTodoByIDFunc = tc.usecaseFunc
			}

			handlers := Handlers{
				UseCase: useCaseMock,
			}

			// act
			handlers.DeleteTodoHandler(rec, req)

			// assert
			if tc.wantCode != rec.Code {
				t.Errorf("unexpected status code: got %d, want %d", rec.Code, tc.wantCode)
			}

			if tc.wantContentType != "" {
				gotCT := rec.Header().Get("Content-Type")
				if gotCT != tc.wantContentType {
					t.Errorf("unexpected Content-Type: got %q, want %q", gotCT, tc.wantContentType)
				}
			}

			if tc.wantBody != "" {
				gotBody := rec.Body.String()
				if gotBody != tc.wantBody {
					t.Errorf("unexpected body: got %q, want %q", gotBody, tc.wantBody)
				}
			}

			if useCaseMock.DeleteTodoByIDCalls != tc.wantCalls {
				t.Errorf("unexpected calls: got %d, want %d", useCaseMock.DeleteTodoByIDCalls, tc.wantCalls)
			}

			if tc.wantCalls > 0 {
				if useCaseMock.LastGetID != tc.wantID {
					t.Errorf("unexpected id: got %d, want %d", useCaseMock.LastGetID, tc.wantID)
				}
			}
		})
	}
}
