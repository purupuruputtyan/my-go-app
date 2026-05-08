package memory

import (
	"errors"
	"github.com/google/uuid"
	"my-go-app/internal/domain/todo"
)

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

func (r *TodoMemory) Create(todo domain.Todo) domain.Todo {
	todo.ID = uuid.NewString()

	r.todos = append(r.todos, todo)

	return todo
}

func (r *TodoMemory) Show(id string) (domain.Todo, error) {
	for _, todo := range r.todos {
		if todo.ID == id {
			return todo, nil
		}
	}

	return domain.Todo{}, errors.New("not found")
}
