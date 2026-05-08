package todo

import "my-go-app/internal/domain/todo"

type TodoUseCase struct {
	repo domain.TodoRepository
}

func NewTodoUseCase(repo domain.TodoRepository) *TodoUseCase {
	return &TodoUseCase{
		repo: repo,
	}
}

func (u *TodoUseCase) FindAll() []domain.Todo {
	return u.repo.FindAll()
}

func (u *TodoUseCase) Create(title string) domain.Todo {
	todo := domain.Todo{
		Title:     title,
		Completed: false,
	}

	return u.repo.Create(todo)
}

func (u *TodoUseCase) Show(id string) (domain.Todo, error) {
	return u.repo.Show(id)
}
