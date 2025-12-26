package memory

import (
	"context"
	"slices"
	"sync"

	"github.com/FooxyS/todo-service/internal/domain"
)

type MemoryTodoRepository struct {
	DB     map[int]domain.Todo
	NextID int
	mu     sync.RWMutex
}

func New() *MemoryTodoRepository {
	return &MemoryTodoRepository{
		DB:     map[int]domain.Todo{},
		NextID: 1,
		mu:     sync.RWMutex{},
	}
}

func (m *MemoryTodoRepository) Save(ctx context.Context, todo domain.Todo) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.NextID
	todo.ID = id
	m.DB[id] = todo
	m.NextID++
	return id, nil
}

func (m *MemoryTodoRepository) GetByID(ctx context.Context, id int) (domain.Todo, error) {
	if err := ctx.Err(); err != nil {
		return domain.Todo{}, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.DB[id]
	if !ok {
		return domain.Todo{}, domain.ErrTodoNotExist
	}

	return v, nil
}

func (m *MemoryTodoRepository) UpdateByID(ctx context.Context, id int, todo domain.Todo) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.DB[id]
	if !ok {
		return domain.ErrTodoNotExist
	}
	todo.ID = id
	m.DB[id] = todo

	return nil
}

func (m *MemoryTodoRepository) DeleteByID(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.DB[id]
	if !ok {
		return domain.ErrTodoNotExist
	}
	delete(m.DB, id)

	return nil
}

func (m *MemoryTodoRepository) ReadAll(ctx context.Context) ([]domain.Todo, error) {
	if err := ctx.Err(); err != nil {
		return []domain.Todo{}, err
	}

	m.mu.RLock()
	res := make([]domain.Todo, 0, len(m.DB))
	for _, v := range m.DB {
		res = append(res, v)
	}
	m.mu.RUnlock()

	slices.SortFunc(res, func(a domain.Todo, b domain.Todo) int {
		if a.ID < b.ID {
			return -1
		} else if a.ID > b.ID {
			return 1
		}
		return 0
	})

	return res, nil
}
