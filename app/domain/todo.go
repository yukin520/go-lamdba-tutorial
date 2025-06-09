package domain

import (
	"context"
	"time"
)

type ToDo struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Completed   bool      `json:"completed"`
}

type TodoRepository interface {
	ListTodo(ctx context.Context) ([]*ToDo, error)
	GetTodo(ctx context.Context, id uint) (*ToDo, error)
	CreateTodo(ctx context.Context, param *ToDo) (uint, error)
	UpdateTodo(ctx context.Context, param *ToDo) (*ToDo, error)
	DeleteTodo(ctx context.Context, id uint) error
}

type TodoUsecase interface {
	ListTodo(ctx context.Context) ([]*ToDo, error)
	GetTodo(ctx context.Context, id uint) (*ToDo, error)
	CreateTodo(ctx context.Context, param *ToDo) (uint, error)
	UpdateTodo(ctx context.Context, param *ToDo) (*ToDo, error)
	DeleteTodo(ctx context.Context, id uint) error
}
