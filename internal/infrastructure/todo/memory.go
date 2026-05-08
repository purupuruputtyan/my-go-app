package memory

import (
	"errors"
	"github.com/google/uuid"
	"my-go-app/internal/domain/todo"
)

type TodoMemory struct {
	todos []todo.Todo
}

func NewTodoMemory() *TodoMemory {
	return &TodoMemory{
		todos: []todo.Todo{},
	}
}

func (r *TodoMemory) FindAll() []todo.Todo {
	return r.todos
}

func (r *TodoMemory) Create(todo todo.Todo) todo.Todo {
	todo.ID = uuid.NewString()

	r.todos = append(r.todos, todo)

	return todo
}

func (r *TodoMemory) Show(id string) (todo.Todo, error) {
	for _, todo := range r.todos {
		if todo.ID == id {
			return todo, nil
		}
	}

	return todo.Todo{}, errors.New("not found")
}
