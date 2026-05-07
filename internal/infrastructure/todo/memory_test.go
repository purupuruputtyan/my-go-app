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
