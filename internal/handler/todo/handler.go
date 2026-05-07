package handler

import (
	"encoding/json"
	"net/http"

	"my-go-app/internal/usecase/todo"
)

type TodoHandler struct {
	usecase *todo.TodoUseCase
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
