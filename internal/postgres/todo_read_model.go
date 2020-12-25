package postgres

import (
	"context"
	"fmt"

	"github.com/go-event-store/pg"
	"github.com/jackc/pgx/v4/pgxpool"
)

const TodoTable = "app_todo"

type TodoReadModel struct {
	client *pg.Client
	stack  []struct {
		method string
		args   []map[string]interface{}
	}
}

func (r *TodoReadModel) Init(ctx context.Context) error {
	_, err := r.client.Conn().(*pgxpool.Pool).Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id UUID NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			done BOOLEAN DEFAULT FALSE,
			deadline TIMESTAMP WITHOUT TIME ZONE NOT NULL,
			updated TIMESTAMP WITHOUT TIME ZONE NOT NULL,
			created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
			PRIMARY KEY (id)
		)`, TodoTable))

	return err
}

func (r *TodoReadModel) IsInitialized(ctx context.Context) (bool, error) {
	return r.client.Exists(ctx, TodoTable)
}

func (r *TodoReadModel) Reset(ctx context.Context) error {
	return r.client.Reset(ctx, TodoTable)
}

func (r *TodoReadModel) Delete(ctx context.Context) error {
	return r.client.Delete(ctx, TodoTable)
}

func (r *TodoReadModel) Stack(method string, args ...map[string]interface{}) {
	r.stack = append(r.stack, struct {
		method string
		args   []map[string]interface{}
	}{method: method, args: args})
}

func (r *TodoReadModel) Persist(ctx context.Context) error {
	var err error
	for _, command := range r.stack {
		switch command.method {
		case "insert":
			err = r.client.Insert(ctx, TodoTable, command.args[0])
			if err != nil {
				return err
			}
		case "remove":
			err = r.client.Remove(ctx, TodoTable, command.args[0])
			if err != nil {
				return err
			}
		case "update":
			err = r.client.Update(ctx, TodoTable, command.args[0], command.args[1])
			if err != nil {
				return err
			}
		}
	}

	r.stack = make([]struct {
		method string
		args   []map[string]interface{}
	}, 0)

	return err
}

func NewTodoReadModel(client *pg.Client) *TodoReadModel {
	return &TodoReadModel{client: client}
}
