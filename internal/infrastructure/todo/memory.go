package memory

import (
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

	return domain.Todo{}, domain.ErrTodoNotFound
}

func (r *TodoMemory) Update(todo domain.Todo) (domain.Todo, error) {
	for i, t := range r.todos {
		if t.ID == todo.ID {
			r.todos[i].Title = todo.Title
			r.todos[i].Completed = todo.Completed

			return r.todos[i], nil
		}
	}

	return domain.Todo{}, domain.ErrTodoNotFound
}
