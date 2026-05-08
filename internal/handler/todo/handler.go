package handler

import (
	"encoding/json"
	"net/http"

	domain "my-go-app/internal/domain/todo"
	"my-go-app/internal/usecase/todo"
)

type TodoHandler struct {
	usecase *todo.TodoUseCase
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func New(usecase *todo.TodoUseCase) *TodoHandler {
	return &TodoHandler{
		usecase: usecase,
	}
}

func (h *TodoHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	todos := h.usecase.FindAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	todo, err := h.usecase.Create(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) Show(w http.ResponseWriter, r *http.Request, id string) {
	todo, err := h.usecase.Show(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	var req UpdateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	todo := domain.Todo{
		ID:        id,
		Title:     req.Title,
		Completed: req.Completed,
	}

	updated, err := h.usecase.Update(todo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	todo, err := h.usecase.Delete(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}
