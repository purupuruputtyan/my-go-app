package handler

import (
	"encoding/json"
	domain "my-go-app/internal/domain/todo"
	memory "my-go-app/internal/infrastructure/todo"
	uc "my-go-app/internal/usecase/todo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTodoHandler_FindAll(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	_ = repo.Create(domain.Todo{Title: "first"})

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	h.FindAll(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got []domain.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(got))
	}
	if got[0].Title != "first" {
		t.Fatalf("expected title first, got %s", got[0].Title)
	}
}

func TestTodoHandler_Create(t *testing.T) {

	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	reqBody := strings.NewReader(`{"title":"first"}`)

	req := httptest.NewRequest(http.MethodPost, "/todos", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var got domain.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Title != "first" {
		t.Fatalf("expected title first, got %s", got.Title)
	}

	todos := repo.FindAll()
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo in repo, got %d", len(todos))
	}
}

func TestTodoHandler_Show(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	created := repo.Create(domain.Todo{
		Title: "first",
	})

	req := httptest.NewRequest(http.MethodGet, "/todos/"+created.ID, nil)
	w := httptest.NewRecorder()

	h.Show(w, req, created.ID)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got domain.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, got.ID)
	}

	if got.Title != "first" {
		t.Fatalf("expected title first, got %s", got.Title)
	}
}
