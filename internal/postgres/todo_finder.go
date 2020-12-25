package postgres

import (
	"context"
	"fmt"

	todo "github.com/go-event-store/example/internal"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TodoFinder struct {
	db *pgxpool.Pool
}

func (f *TodoFinder) Find(ctx context.Context, todoID string) (todo.TodoView, error) {
	view := todo.TodoView{}

	row := f.db.QueryRow(ctx, fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, TodoTable), todoID)
	err := row.Scan(&view.ID, &view.Title, &view.Description, &view.Done, &view.Deadline, &view.UpdatedAt, &view.CreatedAt)
	if err == pgx.ErrNoRows {
		return view, todo.ErrTodoNotFound
	}

	return view, err
}

func (f *TodoFinder) FindAll(ctx context.Context) ([]todo.TodoView, error) {
	list := make([]todo.TodoView, 0)

	rows, err := f.db.Query(ctx, fmt.Sprintf(`SELECT * FROM %s`, TodoTable))
	if err != nil {
		return list, err
	}

	for rows.Next() {
		view := todo.TodoView{}

		err := rows.Scan(&view.ID, &view.Title, &view.Description, &view.Done, &view.Deadline, &view.UpdatedAt, &view.CreatedAt)
		if err != nil {
			return make([]todo.TodoView, 0), err
		}

		list = append(list, view)
	}

	return list, nil
}

func NewTodoFinder(db *pgxpool.Pool) *TodoFinder {
	return &TodoFinder{db}
}
