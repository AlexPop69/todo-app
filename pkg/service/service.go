package service

import "github.com/AlexPop69/todo-app/pkg/repository"

type Autorization interface {
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Autorization
	TodoList
	TodoItem
}

func NewService(rep *repository.Repository) *Service {
	return &Service{}
}
