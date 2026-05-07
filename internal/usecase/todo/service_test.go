package todo

import (
	domain "my-go-app/internal/domain/todo"
	"testing"
)

type stubRepo struct {
	todos []domain.Todo
}

func (s *stubRepo) FindAll() []domain.Todo {
	return s.todos
}

func (s *stubRepo) Create(t domain.Todo) domain.Todo {
	s.todos = append(s.todos, t)
	return t
}

func TestTodoUseCase_FindAll(t *testing.T) {
	repo := &stubRepo{
		todos: []domain.Todo{{ID: "1", Title: "learn go", Completed: false}},
	}
	uc := NewTodoUseCase(repo)

	todos := uc.FindAll()
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	if todos[0].Title != "learn go" {
		t.Fatalf("expected title learn go, got %s", todos[0].Title)
	}
}
