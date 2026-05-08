package memory

import (
	"my-go-app/internal/domain/todo"
	"testing"
)

func TestTodoMemory_FindAll(t *testing.T) {
	repo := NewTodoMemory()

	_ = repo.Create(todo.Todo{Title: "first"})

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

	todo := todo.Todo{
		Title: "first",
	}

	created := repo.Create(todo)

	if created.Title != "first" {
		t.Fatalf("expected title first, got %s", created.Title)
	}

	todos := repo.FindAll()

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}

	if todos[0].Title != "first" {
		t.Fatalf("expected title first, got %s", todos[0].Title)
	}
}

func TestTodoMemory_Show(t *testing.T) {
	repo := NewTodoMemory()

	repo.todos = []todo.Todo{
		{ID: "111", Title: "テスト", Completed: true},
	}

	result, err := repo.Show("111")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "111" {
		t.Fatalf("expected id 111, got %s", result.ID)
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

	repo.todos = []todo.Todo{
		{ID: "111", Title: "テスト", Completed: true},
	}

	_, err := repo.Show("999")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
