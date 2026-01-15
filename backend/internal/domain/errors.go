package domain

import "errors"

var (
	ErrNoTitle      = errors.New("title is empty")
	ErrTodoNotExist = errors.New("todo with specified id does not exist")
)
