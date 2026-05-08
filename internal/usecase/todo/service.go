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

func (u *TodoUseCase) Create(title string) (domain.Todo, error) {
	if title == "" {
		return domain.Todo{}, domain.ErrTitleRequired
	}
	if len(title) > 100 {
		return domain.Todo{}, domain.ErrTitleTooLong
	}

	todo := domain.Todo{
		Title:     title,
		Completed: false,
	}

	return u.repo.Create(todo), nil
}

func (u *TodoUseCase) Show(id string) (domain.Todo, error) {
	return u.repo.Show(id)
}

func (u *TodoUseCase) Update(todo domain.Todo) (domain.Todo, error) {
	if todo.Title == "" {
		return domain.Todo{}, domain.ErrTitleRequired
	}

	if len(todo.Title) > 100 {
		return domain.Todo{}, domain.ErrTitleTooLong
	}

	return u.repo.Update(todo)
}
