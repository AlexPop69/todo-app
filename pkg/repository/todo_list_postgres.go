package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlexPop69/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Add(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	addListQuery := fmt.Sprintf(`INSERT INTO %s (title, description) 
		VALUES ($1, $2) RETURNING id`, todoListsTable)
	row := tx.QueryRow(addListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	addUsersListQuery := fmt.Sprintf(`INSERT INTO %s (user_id, list_id)
		VALUES ($1, $2)`, usersListsTable)
	_, err = tx.Exec(addUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s
	 	tl INNER JOIN %s ul on tl.id = ul.list_id 
	 	WHERE ul.user_id=$1`, todoListsTable, usersListsTable)

	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s
	tl INNER JOIN %s ul on tl.id = ul.list_id 
	WHERE ul.user_id=$1 AND ul.list_id=$2`, todoListsTable, usersListsTable)

	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf(`DELETE FROM %s tl	
		USING %s ul
		WHERE tl.id=ul.list_id AND ul.user_id=$1 AND ul.list_id=$2`,
		todoListsTable, usersListsTable)

	del, err := r.db.Exec(query, userId, listId)

	result, err := del.RowsAffected()
	if result == 0 {
		err = errors.New("list does not exist")
	}

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *&input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *&input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d`,
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args...)

	_, err := r.db.Exec(query, args...)

	return err
}
