package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"my-go-app/internal/domain/todo"
)

func TestTodoPostgres_CreateAndShow(t *testing.T) {
	db := openTestDB(t)
	repo := NewTodoPostgres(db)

	created := repo.Create(domain.Todo{
		Title:     "learn go",
		Completed: false,
	})
	t.Cleanup(func() {
		cleanupTodo(t, db, created.ID)
	})

	if created.ID == "" {
		t.Fatalf("expected id to be set")
	}

	if created.Title != "learn go" {
		t.Fatalf("expected title learn go, got %s", created.Title)
	}

	if created.Completed {
		t.Fatalf("expected completed false")
	}

	found, err := repo.Show(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, found.ID)
	}

	if found.Title != "learn go" {
		t.Fatalf("expected title learn go, got %s", found.Title)
	}
}

func TestTodoPostgres_FindAll(t *testing.T) {
	db := openTestDB(t)
	repo := NewTodoPostgres(db)

	first := repo.Create(domain.Todo{Title: "first"})
	second := repo.Create(domain.Todo{Title: "second"})

	t.Cleanup(func() {
		cleanupTodo(t, db, first.ID)
		cleanupTodo(t, db, second.ID)
	})

	todos := repo.FindAll()

	if !containsTodoID(todos, first.ID) {
		t.Fatalf("expected todos to contain id %s", first.ID)
	}

	if !containsTodoID(todos, second.ID) {
		t.Fatalf("expected todos to contain id %s", second.ID)
	}
}

func TestTodoPostgres_Update(t *testing.T) {
	db := openTestDB(t)
	repo := NewTodoPostgres(db)

	created := repo.Create(domain.Todo{
		Title:     "before",
		Completed: false,
	})
	t.Cleanup(func() {
		cleanupTodo(t, db, created.ID)
	})

	updated, err := repo.Update(domain.Todo{
		ID:        created.ID,
		Title:     "after",
		Completed: true,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.Title != "after" {
		t.Fatalf("expected title after, got %s", updated.Title)
	}

	if !updated.Completed {
		t.Fatalf("expected completed true")
	}

	found, err := repo.Show(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.Title != "after" {
		t.Fatalf("expected title after, got %s", found.Title)
	}

	if !found.Completed {
		t.Fatalf("expected completed true")
	}
}

func TestTodoPostgres_Update_NotFound(t *testing.T) {
	db := openTestDB(t)
	repo := NewTodoPostgres(db)

	_, err := repo.Update(domain.Todo{
		ID:        "not-found-id",
		Title:     "after",
		Completed: true,
	})

	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestTodoPostgres_Delete(t *testing.T) {
	db := openTestDB(t)
	repo := NewTodoPostgres(db)

	created := repo.Create(domain.Todo{
		Title:     "delete me",
		Completed: false,
	})

	deleted, err := repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if deleted.ID != created.ID {
		t.Fatalf("expected deleted id %s, got %s", created.ID, deleted.ID)
	}

	_, err = repo.Show(created.ID)
	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestTodoPostgres_Delete_NotFound(t *testing.T) {
	db := openTestDB(t)
	repo := NewTodoPostgres(db)

	_, err := repo.Delete("not-found-id")

	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envOrDefault("DB_HOST", "localhost"),
		envOrDefault("DB_PORT", "5432"),
		envOrDefault("DB_USER", "my_go_app"),
		envOrDefault("DB_PASSWORD", "password"),
		envOrDefault("DB_NAME", "my_go_app"),
		envOrDefault("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func cleanupTodo(t *testing.T, db *sql.DB, id string) {
	t.Helper()

	if id == "" {
		return
	}

	if _, err := db.Exec(`DELETE FROM todos WHERE id = $1`, id); err != nil {
		t.Fatalf("failed to cleanup todo %s: %v", id, err)
	}
}

func containsTodoID(todos []domain.Todo, id string) bool {
	for _, todo := range todos {
		if todo.ID == id {
			return true
		}
	}

	return false
}

func envOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
