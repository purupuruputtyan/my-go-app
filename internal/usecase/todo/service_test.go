package todo

import (
	"testing"

	"my-go-app/internal/domain/todo"
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
	return domain.Todo{}, domain.ErrTodoNotFound
}

func (s *stubRepo) Update(t domain.Todo) (domain.Todo, error) {
	for i, todo := range s.todos {
		if todo.ID == t.ID {
			s.todos[i] = t
			return t, nil
		}
	}

	return domain.Todo{}, domain.ErrTodoNotFound
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

	created, err := uc.Create("learn go")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.Title != "learn go" {
		t.Fatalf("expected title learn go, got %s", created.Title)
	}

	if created.Completed {
		t.Fatalf("expected completed false")
	}

	if len(repo.todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(repo.todos))
	}
}

func TestTodoUseCase_Create_EmptyTitle(t *testing.T) {
	repo := &stubRepo{}
	uc := NewTodoUseCase(repo)

	_, err := uc.Create("")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestTodoUseCase_Create_TitleTooLong(t *testing.T) {
	repo := &stubRepo{}
	uc := NewTodoUseCase(repo)

	longTitle := "a"
	for len(longTitle) <= 100 {
		longTitle += "a"
	}

	_, err := uc.Create(longTitle)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrTitleTooLong {
		t.Fatalf("expected ErrTitleTooLong, got %v", err)
	}
}

func TestTodoUseCase_Show(t *testing.T) {
	repo := &stubRepo{
		todos: []domain.Todo{
			{ID: "1", Title: "learn go", Completed: false},
		},
	}
	uc := NewTodoUseCase(repo)

	result, err := uc.Show("1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" {
		t.Fatalf("expected id 1, got %s", result.ID)
	}

	if result.Title != "learn go" {
		t.Fatalf("expected title learn go, got %s", result.Title)
	}
}

func TestTodoUseCase_Show_NotFound(t *testing.T) {
	repo := &stubRepo{}
	uc := NewTodoUseCase(repo)

	_, err := uc.Show("not-found-id")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestTodoUseCase_Update(t *testing.T) {
	repo := &stubRepo{
		todos: []domain.Todo{
			{ID: "1", Title: "before", Completed: false},
		},
	}
	uc := NewTodoUseCase(repo)

	updateTodo := domain.Todo{
		ID:        "1",
		Title:     "after",
		Completed: true,
	}

	updated, err := uc.Update(updateTodo)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.ID != "1" {
		t.Fatalf("expected id 1, got %s", updated.ID)
	}

	if updated.Title != "after" {
		t.Fatalf(
			"expected title after, got %s",
			updated.Title,
		)
	}

	if !updated.Completed {
		t.Fatalf("expected completed true")
	}
}

func TestTodoUseCase_Update_EmptyTitle(t *testing.T) {
	repo := &stubRepo{
		todos: []domain.Todo{
			{ID: "1", Title: "before", Completed: false},
		},
	}
	uc := NewTodoUseCase(repo)

	updateTodo := domain.Todo{
		ID:        "1",
		Title:     "",
		Completed: true,
	}

	_, err := uc.Update(updateTodo)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestTodoUseCase_Update_TitleTooLong(t *testing.T) {
	repo := &stubRepo{
		todos: []domain.Todo{
			{ID: "1", Title: "before", Completed: false},
		},
	}
	uc := NewTodoUseCase(repo)

	longTitle := "a"
	for len(longTitle) <= 100 {
		longTitle += "a"
	}

	updateTodo := domain.Todo{
		ID:        "1",
		Title:     longTitle,
		Completed: true,
	}

	_, err := uc.Update(updateTodo)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrTitleTooLong {
		t.Fatalf("expected ErrTitleTooLong, got %v", err)
	}
}
