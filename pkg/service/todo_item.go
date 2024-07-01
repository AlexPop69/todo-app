package service

import (
	"errors"

	"github.com/AlexPop69/todo-app"
	"github.com/AlexPop69/todo-app/pkg/repository"
)

type TodoItemService struct {
	rep     repository.TodoItem
	listRep repository.TodoList
}

func NewTodoItemService(rep repository.TodoItem, listRep repository.TodoList) *TodoItemService {
	return &TodoItemService{rep: rep, listRep: listRep}
}

func (s *TodoItemService) Add(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRep.GetById(userId, listId)
	if err != nil {
		err = errors.New("list does not exist or does not belongs to user")
		return 0, err
	}

	return s.rep.Add(userId, item)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.rep.GetAll(userId, listId)
}
