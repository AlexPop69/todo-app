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

func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.rep.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.rep.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.rep.Update(userId, itemId, input)
}
