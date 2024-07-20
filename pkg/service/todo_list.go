package service

import (
	"github.com/AlexPop69/todo-app"
	"github.com/AlexPop69/todo-app/pkg/repository"
)

type TodoListService struct {
	rep repository.TodoList
}

func NewTodoListService(rep repository.TodoList) *TodoListService {
	return &TodoListService{rep: rep}
}

func (s *TodoListService) Add(userId int, list todo.TodoList) (int, error) {
	return s.rep.Add(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.rep.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
	return s.rep.GetById(userId, listId)
}

func (s *TodoListService) Delete(userId, listId int) error {
	return s.rep.Delete(userId, listId)
}

func (s *TodoListService) Update(userId, listId int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.rep.Update(userId, listId, input)
}
