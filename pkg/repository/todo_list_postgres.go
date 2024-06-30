package repository

import (
	"errors"
	"fmt"

	"github.com/AlexPop69/todo-app"
	"github.com/jmoiron/sqlx"
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
