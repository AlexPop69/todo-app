package service

import (
	"github.com/AlexPop69/todo-app"
	"github.com/AlexPop69/todo-app/pkg/repository"
)

type TodoListServise struct {
	rep repository.TodoList
}

func NewTodoListService(rep repository.TodoList) *TodoListServise {
	return &TodoListServise{rep: rep}
}

func (s *TodoListServise) Add(userId int, list todo.TodoList) (int, error) {
	return s.rep.Add(userId, list)
}

func (s *TodoListServise) GetAll(userId int) ([]todo.TodoList, error) {
	return s.rep.GetAll(userId)
}

func (s *TodoListServise) GetById(userId, listId int) (todo.TodoList, error) {
	return s.rep.GetById(userId, listId)
}
