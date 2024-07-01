package service

import (
	"github.com/AlexPop69/todo-app"
	"github.com/AlexPop69/todo-app/pkg/repository"
)

type Autorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Add(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todo.UpdateListInput) error
}

type TodoItem interface {
	Add(userId, listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
}

type Service struct {
	Autorization
	TodoList
	TodoItem
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		Autorization: NewAuthService(rep.Autorization),
		TodoList:     NewTodoListService(rep.TodoList),
		TodoItem:     NewTodoItemService(rep.TodoItem, rep.TodoList),
	}
}
