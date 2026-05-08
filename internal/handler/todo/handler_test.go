package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domain "my-go-app/internal/domain/todo"
	memory "my-go-app/internal/infrastructure/todo"
	uc "my-go-app/internal/usecase/todo"
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

func TestTodoHandler_Create_EmptyTitle(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	reqBody := strings.NewReader(`{"title":""}`)

	req := httptest.NewRequest(http.MethodPost, "/todos", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestTodoHandler_Create_TitleTooLong(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	longTitle := "a"
	for len(longTitle) <= 100 {
		longTitle += "a"
	}

	reqBody := strings.NewReader(`{"title":"` + longTitle + `"}`)

	req := httptest.NewRequest(http.MethodPost, "/todos", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestTodoHandler_Create_InvalidJSON(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	reqBody := strings.NewReader(`{"title":`)

	req := httptest.NewRequest(http.MethodPost, "/todos", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	todos := repo.FindAll()
	if len(todos) != 0 {
		t.Fatalf("expected 0 todo in repo, got %d", len(todos))
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

func TestTodoHandler_Show_NotFound(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	req := httptest.NewRequest(http.MethodGet, "/todos/not-found-id", nil)
	w := httptest.NewRecorder()

	h.Show(w, req, "not-found-id")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestTodoHandler_Update(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	created := repo.Create(domain.Todo{
		Title:     "before",
		Completed: false,
	})

	reqBody := strings.NewReader(`{
		"title":"after",
		"completed":true
	}`)

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/"+created.ID,
		reqBody,
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Update(w, req, created.ID)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got domain.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.ID != created.ID {
		t.Fatalf(
			"expected id %s, got %s",
			created.ID,
			got.ID,
		)
	}

	if got.Title != "after" {
		t.Fatalf(
			"expected title after, got %s",
			got.Title,
		)
	}

	if !got.Completed {
		t.Fatalf("expected completed true")
	}
}

func TestTodoHandler_Update_EmptyTitle(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	created := repo.Create(domain.Todo{
		Title:     "before",
		Completed: false,
	})

	reqBody := strings.NewReader(`{
		"title":"",
		"completed":true
	}`)

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/"+created.ID,
		reqBody,
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Update(w, req, created.ID)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestTodoHandler_Update_TitleTooLong(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	created := repo.Create(domain.Todo{
		Title:     "before",
		Completed: false,
	})

	longTitle := "a"
	for len(longTitle) <= 100 {
		longTitle += "a"
	}

	reqBody := strings.NewReader(`{"title":"` + longTitle + `"}`)

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/"+created.ID,
		reqBody,
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Update(w, req, created.ID)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestTodoHandler_Update_InvalidJSON(t *testing.T) {
	repo := memory.NewTodoMemory()
	usecase := uc.NewTodoUseCase(repo)
	h := New(usecase)

	created := repo.Create(domain.Todo{
		Title:     "before",
		Completed: false,
	})

	reqBody := strings.NewReader(`{"title":`)

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/"+created.ID,
		reqBody,
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Update(w, req, created.ID)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	todos := repo.FindAll()
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo in repo, got %d", len(todos))
	}

	if todos[0].Title != "before" {
		t.Fatalf(
			"expected title before, got %s",
			todos[0].Title,
		)
	}

	if todos[0].Completed {
		t.Fatalf("expected completed false")
	}
}
