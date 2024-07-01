package repository

import (
	"fmt"

	"github.com/AlexPop69/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Add(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int

	addItemQuery := fmt.Sprintf(`INSERT INTO %s (title, description) 
		VALUES ($1, $2) RETURNING id`, todoItemsTable)

	row := tx.QueryRow(addItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	addListItemQuery := fmt.Sprintf(`INSERT INTO %s (list_id, item_id)
										VALUES ($1, $2)`, listsItemsTable)
	logrus.Println("addListItemQuery:", listId, itemId)
	_, err = tx.Exec(addListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}
	logrus.Println("end")
	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done
			FROM %s ti INNER JOIN %s li on li.item_id = ti.id
			INNER JOIN %s ul on ul.list_id = li.list_id 
	 		WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	logrus.Println("GetAll:", listId, userId)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}
