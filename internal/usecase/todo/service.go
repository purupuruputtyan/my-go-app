package todo

import "my-go-app/internal/domain/todo"

type TodoUseCase struct {
	repo todo.TodoRepository
}

func NewTodoUseCase(repo todo.TodoRepository) *TodoUseCase {
	return &TodoUseCase{
		repo: repo,
	}
}

func (u *TodoUseCase) FindAll() []todo.Todo {
	return u.repo.FindAll()
}

func (u *TodoUseCase) Create(title string) todo.Todo {
	todo := todo.Todo{
		Title:     title,
		Completed: false,
	}

	return u.repo.Create(todo)
}

func (u *TodoUseCase) Show(id string) (todo.Todo, error) {
	return u.repo.Show(id)
}
