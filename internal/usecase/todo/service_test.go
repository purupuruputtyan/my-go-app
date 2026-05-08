package todo

import (
	"errors"
	"my-go-app/internal/domain/todo"
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

func (s *stubRepo) Show(id string) (domain.Todo, error) {
	for _, t := range s.todos {
		if t.ID == id {
			return t, nil
		}
	}
	return domain.Todo{}, errors.New("not found")
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

func TestTodoUseCase_Create(t *testing.T) {
	repo := &stubRepo{}
	uc := NewTodoUseCase(repo)

	created := uc.Create("learn go")
	if created.Title != "learn go" {
		t.Fatalf("expected title learn go, got %s", created.Title)
	}

	if len(repo.todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(repo.todos))
	}
}
