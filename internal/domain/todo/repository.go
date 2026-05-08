package domain

type TodoRepository interface {
	FindAll() []Todo
	Create(todo Todo) Todo
	Show(id string) (Todo, error)
	Update(todo Todo) (Todo, error)
	Delete(id string) (Todo, error)
}
