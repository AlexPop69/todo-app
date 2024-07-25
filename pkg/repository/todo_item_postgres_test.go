package repository

import (
	"errors"

	"testing"

	"github.com/AlexPop69/todo-app"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoItemPostgres_Add(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item   todo.TodoItem
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		input        args
		id           int
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Ok",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(`INSERT INTO todo_items`).
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				mock.ExpectExec(`INSERT INTO lists_items`).
					WithArgs(args.listId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Rollback first insert",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("insert error"))
				mock.ExpectQuery(`INSERT INTO todo_items`).
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Rollback second insert",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(`INSERT INTO todo_items`).
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				mock.ExpectExec(`INSERT INTO lists_items`).
					WithArgs(args.listId, id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.id)

			actualId, err := r.Add(testCase.input.listId, testCase.input.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, actualId)
			}

		})
	}

}

func TestTodoItemPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		userId int
	}

	testTable := []struct {
		name         string
		input        args
		mockBehavior func()
		want         []todo.TodoItem
	}{
		{
			name: "Ok",
			input: args{
				listId: 1,
				userId: 1,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "test title 1", "test description 1", true).
					AddRow(2, "test title 2", "test description 2", false).
					AddRow(3, "test title 3", "test description 3", false)

				mock.ExpectQuery(`SELECT (.+)
					FROM todo_items ti INNER JOIN lists_items li on (.+)
					INNER JOIN users_lists ul on (.+)
					WHERE (.+)`).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			want: []todo.TodoItem{
				{1, "test title 1", "test description 1", true},
				{2, "test title 2", "test description 2", false},
				{3, "test title 3", "test description 3", false},
			},
		},
		{
			name: "No records",
			input: args{
				listId: 1,
				userId: 1,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery(`SELECT (.+)
					FROM todo_items ti INNER JOIN lists_items li on (.+)
					INNER JOIN users_lists ul on (.+)
					WHERE (.+)`).
					WithArgs(1, 1).WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			actualResult, err := r.GetAll(testCase.input.userId, testCase.input.listId)

			assert.NoError(t, err)
			assert.Equal(t, testCase.want, actualResult)

		})
	}

}

func TestTodoItemPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		userId int
		itemId int
	}

	testTable := []struct {
		name         string
		input        args
		mockBehavior func()
		want         todo.TodoItem
		wantErr      bool
	}{
		{
			name: "Ok",
			input: args{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "test title", "test description", true)

				mock.ExpectQuery(`SELECT (.+)
					FROM todo_items ti INNER JOIN lists_items li on (.+)
					INNER JOIN users_lists ul on (.+)
					WHERE (.+)`).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			want:    todo.TodoItem{1, "test title", "test description", true},
			wantErr: false,
		},
		{
			name: "No record",
			input: args{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery(`SELECT (.+)
					FROM todo_items ti INNER JOIN lists_items li on (.+)
					INNER JOIN users_lists ul on (.+)
					WHERE (.+)`).
					WithArgs(2, 2).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			actualResult, err := r.GetById(testCase.input.userId, testCase.input.itemId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, actualResult)
			}

		})
	}

}

func TestTodoItemPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		userId int
		itemId int
	}

	testTable := []struct {
		name         string
		input        args
		mockBehavior func()
		wantErr      bool
	}{
		{
			name: "Ok",
			input: args{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func() {
				mock.ExpectExec(`DELETE FROM todo_items ti 
					USING lists_items li, users_lists ul 
					WHERE (.+) `).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "No record",
			input: args{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func() {
				mock.ExpectExec(`DELETE FROM todo_items ti 
					USING lists_items li, users_lists ul 
					WHERE (.+) `).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err := r.Delete(testCase.input.userId, testCase.input.itemId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTodoItemPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		userId int
		itemId int
		input  todo.UpdateItemInput
	}

	testTable := []struct {
		name         string
		input        args
		mockBehavior func()
		wantErr      bool
	}{
		{
			name: "Ok",
			input: args{
				userId: 1,
				itemId: 1,
				input: todo.UpdateItemInput{
					Title:       stringPointer("test title"),
					Description: stringPointer("test description"),
					Done:        boolPointer(true),
				},
			},
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE todo_items ti SET (.+)
					FROM lists_items li, users_lists ul
					WHERE (.+)`).
					WillReturnResult(sqlmock.NewResult(0, 1))

			},
			wantErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err := r.Update(testCase.input.userId, testCase.input.itemId, testCase.input.input)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
