package repository

import (
	"errors"
	"testing"

	"github.com/AlexPop69/todo-app"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoListPostgres_Add(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		list   todo.TodoList
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		input        args
		wantId       int
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Ok",
			input: args{
				userId: 1,
				list: todo.TodoList{
					Title:       "test title",
					Description: "test description",
				},
			},
			wantId: 1,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`INSERT INTO todo_lists`).
					WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec(`INSERT INTO users_lists`).
					WithArgs(args.userId, id).
					WillReturnResult(sqlmock.NewResult(int64(id), 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Rollback first insert",
			input: args{
				userId: 1,
				list: todo.TodoList{
					Title:       "",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("list insert error"))
				mock.ExpectQuery(`INSERT INTO todo_lists`).
					WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Rollback second insert",
			input: args{
				userId: 1,
				list: todo.TodoList{
					Title:       "test title",
					Description: "test description",
				},
			},
			wantId: 1,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(`INSERT INTO todo_lists`).
					WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec(`INSERT INTO users_lists`).
					WithArgs(args.userId, id).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.wantId)

			actualId, err := r.Add(testCase.input.userId, testCase.input.list)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.wantId, actualId)
			}

		})
	}

}

func TestTodoListostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
	}

	testTable := []struct {
		name         string
		input        args
		mockBehavior func()
		want         []todo.TodoList
	}{
		{
			name: "Ok",
			input: args{
				userId: 1,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "test title 1", "test description 1").
					AddRow(2, "test title 2", "test description 2").
					AddRow(3, "test title 3", "test description 3")

				mock.ExpectQuery(`SELECT (.+)
					FROM todo_lists tl
					INNER JOIN users_lists ul on (.+)
					WHERE (.+)`).
					WithArgs(1).WillReturnRows(rows)
			},
			want: []todo.TodoList{
				{1, "test title 1", "test description 1"},
				{2, "test title 2", "test description 2"},
				{3, "test title 3", "test description 3"},
			},
		},
		{
			name: "No records",
			input: args{
				userId: 1,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})

				mock.ExpectQuery(`SELECT (.+)
				FROM todo_lists tl
				INNER JOIN users_lists ul on (.+)
				WHERE (.+)`).
					WithArgs(1).WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			actualResult, err := r.GetAll(testCase.input.userId)

			assert.NoError(t, err)
			assert.Equal(t, testCase.want, actualResult)

		})
	}

}

func TestTodoListPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		listId int
	}

	testTable := []struct {
		name         string
		input        args
		mockBehavior func()
		want         todo.TodoList
		wantErr      bool
	}{
		{
			name: "Ok",
			input: args{
				userId: 1,
				listId: 1,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "test title", "test description")

				mock.ExpectQuery(`SELECT (.+)
					FROM todo_lists tl
					INNER JOIN users_lists ul on (.+)
					WHERE (.+)`).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			want:    todo.TodoList{1, "test title", "test description"},
			wantErr: false,
		},
		{
			name: "No record",
			input: args{
				userId: 1,
				listId: 1,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})

				mock.ExpectQuery(`SELECT (.+)
				FROM todo_lists tl
				INNER JOIN users_lists ul on (.+)
				WHERE (.+)`).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			actualResult, err := r.GetById(testCase.input.userId, testCase.input.listId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, actualResult)
			}

		})
	}

}

func TestTodoListPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		listId int
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
				listId: 1,
			},
			mockBehavior: func() {
				mock.ExpectExec(`DELETE FROM todo_lists tl
					USING users_lists ul
					WHERE (.+) `).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "No record",
			input: args{
				userId: 1,
				listId: 1,
			},
			mockBehavior: func() {
				mock.ExpectExec(`DELETE FROM todo_lists tl
					USING users_lists ul
					WHERE (.+) `).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err := r.Delete(testCase.input.userId, testCase.input.listId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTodoListPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		userId int
		listId int
		input  todo.UpdateListInput
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
				listId: 1,
				input: todo.UpdateListInput{
					Title:       stringPointer("test title"),
					Description: stringPointer("test description"),
				},
			},
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE todo_lists tl
					SET (.+)
					FROM users_lists ul
					WHERE (.+)`).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err := r.Update(testCase.input.userId, testCase.input.listId, testCase.input.input)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
