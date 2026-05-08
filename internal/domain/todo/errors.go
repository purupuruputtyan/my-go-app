package domain

import "errors"

var (
	ErrTodoNotFound  = errors.New("todo not found")
	ErrTitleRequired = errors.New("title is required")
	ErrTitleTooLong  = errors.New("title too long (max 100 characters)")
)
