package handler

import (
	"encoding/json"
	"net/http"

	"my-go-app/internal/usecase/todo"
)

type TodoHandler struct {
	usecase *todo.TodoUseCase
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func New(usecase *todo.TodoUseCase) *TodoHandler {
	return &TodoHandler{
		usecase: usecase,
	}
}

func (h *TodoHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	todos := h.usecase.FindAll()
	_ = json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	todo := h.usecase.Create(req.Title)

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

	_ = json.NewEncoder(w).Encode(todo)
}
