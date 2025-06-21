package usecase

import (
	"context"

	"github.com/yukin520/go-lamdba-tutorial/app/domain"
)

type usecase struct {
	todoRepo domain.TodoRepository
}

func NewUsecase(r domain.TodoRepository) domain.TodoUsecase {
	return &usecase{
		todoRepo: r,
	}
}

func (m *usecase) ListTodo(ctx context.Context) ([]*domain.ToDo, error) {
	res, err := m.todoRepo.ListTodo(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *usecase) GetTodo(ctx context.Context, id uint) (*domain.ToDo, error) {
	res, err := m.todoRepo.GetTodo(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *usecase) CreateTodo(ctx context.Context, param *domain.ToDo) (uint, error) {
	panic("GetTodo not implemented")
}
func (m *usecase) UpdateTodo(ctx context.Context, param *domain.ToDo) (*domain.ToDo, error) {
	panic("GetTodo not implemented")
}
func (m *usecase) DeleteTodo(ctx context.Context, id uint) error {
	panic("GetTodo not implemented")
}
