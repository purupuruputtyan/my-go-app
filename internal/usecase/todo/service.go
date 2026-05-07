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
