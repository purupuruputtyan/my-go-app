package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"

	"my-go-app/internal/domain/todo"
	"my-go-app/internal/infrastructure/sqlboiler/models"
)

type TodoPostgres struct {
	db *sql.DB
}

func NewTodoPostgres(db *sql.DB) *TodoPostgres {
	return &TodoPostgres{db: db}
}

func (r *TodoPostgres) FindAll() []domain.Todo {
	ctx := context.Background()

	rows, err := models.Todos(
		qm.OrderBy(models.TodoColumns.CreatedAt+" DESC"),
	).All(ctx, r.db)
	if err != nil {
		return []domain.Todo{}
	}

	todos := make([]domain.Todo, 0, len(rows))
	for _, row := range rows {
		todos = append(todos, toDomainTodo(row))
	}

	return todos
}

func (r *TodoPostgres) Create(input domain.Todo) domain.Todo {
	ctx := context.Background()

	row := &models.Todo{
		ID:        uuid.NewString(),
		Title:     input.Title,
		Completed: input.Completed,
	}

	if err := row.Insert(ctx, r.db, boil.Infer()); err != nil {
		return domain.Todo{}
	}

	return toDomainTodo(row)
}

func (r *TodoPostgres) Show(id string) (domain.Todo, error) {
	ctx := context.Background()

	row, err := models.FindTodo(ctx, r.db, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Todo{}, domain.ErrTodoNotFound
	}
	if err != nil {
		return domain.Todo{}, err
	}

	return toDomainTodo(row), nil
}

func (r *TodoPostgres) Update(input domain.Todo) (domain.Todo, error) {
	ctx := context.Background()

	row, err := models.FindTodo(ctx, r.db, input.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Todo{}, domain.ErrTodoNotFound
	}
	if err != nil {
		return domain.Todo{}, err
	}

	row.Title = input.Title
	row.Completed = input.Completed

	rowsAffected, err := row.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return domain.Todo{}, err
	}
	if rowsAffected == 0 {
		return domain.Todo{}, domain.ErrTodoNotFound
	}

	return toDomainTodo(row), nil
}

func (r *TodoPostgres) Delete(id string) (domain.Todo, error) {
	ctx := context.Background()

	row, err := models.FindTodo(ctx, r.db, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Todo{}, domain.ErrTodoNotFound
	}
	if err != nil {
		return domain.Todo{}, err
	}

	deleted := toDomainTodo(row)

	rowsAffected, err := row.Delete(ctx, r.db)
	if err != nil {
		return domain.Todo{}, err
	}
	if rowsAffected == 0 {
		return domain.Todo{}, domain.ErrTodoNotFound
	}

	return deleted, nil
}

func toDomainTodo(row *models.Todo) domain.Todo {
	return domain.Todo{
		ID:        row.ID,
		Title:     row.Title,
		Completed: row.Completed,
	}
}
