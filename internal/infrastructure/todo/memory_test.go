package memory

import (
	"testing"

	"my-go-app/internal/domain/todo"
)

func TestTodoMemory_FindAll(t *testing.T) {
	repo := NewTodoMemory()

	_ = repo.Create(domain.Todo{Title: "first"})

	todos := repo.FindAll()
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	if todos[0].Title != "first" {
		t.Fatalf("expected first todo title to be first, got %s", todos[0].Title)
	}
}

func TestTodoMemory_Create(t *testing.T) {
	repo := NewTodoMemory()

	todo := domain.Todo{
		Title: "first",
	}

	created := repo.Create(todo)

	if created.ID == "" {
		t.Fatalf("expected id to be set")
	}

	if created.Title != "first" {
		t.Fatalf("expected title first, got %s", created.Title)
	}

	todos := repo.FindAll()

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}

	if todos[0].ID == "" {
		t.Fatalf("expected id to be set")
	}

	if todos[0].Title != "first" {
		t.Fatalf("expected title first, got %s", todos[0].Title)
	}
}

func TestTodoMemory_Show(t *testing.T) {
	repo := NewTodoMemory()

	created := repo.Create(domain.Todo{
		Title:     "テスト",
		Completed: true,
	})

	result, err := repo.Show(created.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, result.ID)
	}

	if result.Title != "テスト" {
		t.Fatalf("expected title テスト, got %s", result.Title)
	}

	if !result.Completed {
		t.Fatalf("expected completed true")
	}
}

func TestTodoMemory_Show_NotFound(t *testing.T) {
	repo := NewTodoMemory()

	_ = repo.Create(domain.Todo{
		Title:     "テスト",
		Completed: true,
	})

	_, err := repo.Show("not-found-id")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestTodoMemory_Update(t *testing.T) {
	repo := NewTodoMemory()

	todo := domain.Todo{
		Title: "first",
	}

	created := repo.Create(todo)

	updateTodo := domain.Todo{
		ID:        created.ID,
		Title:     "更新テスト",
		Completed: true,
	}

	updated, err := repo.Update(updateTodo)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.ID != created.ID {
		t.Fatalf(
			"expected id %s, got %s",
			created.ID,
			updated.ID,
		)
	}

	if updated.Title != "更新テスト" {
		t.Fatalf(
			"expected title 更新テスト, got %s",
			updated.Title,
		)
	}

	if !updated.Completed {
		t.Fatalf("expected completed true")
	}

	todos := repo.FindAll()

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}

	if todos[0].Title != "更新テスト" {
		t.Fatalf(
			"expected title 更新テスト, got %s",
			todos[0].Title,
		)
	}

	if !todos[0].Completed {
		t.Fatalf("expected completed true")
	}
}
