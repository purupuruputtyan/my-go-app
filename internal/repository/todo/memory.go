package repository

import "my-go-app/internal/domain/todo"

type TodoMemory struct {
	todos []domain.Todo
}

func NewTodoMemory() *TodoMemory {
	return &TodoMemory{
		todos: []domain.Todo{},
	}
}

func (r *TodoMemory) FindAll() []domain.Todo {
	return r.todos
}
