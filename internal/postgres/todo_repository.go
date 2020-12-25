package postgres

import (
	"context"

	eventstore "github.com/go-event-store/eventstore"
	todo "github.com/go-event-store/example/internal"
	uuid "github.com/satori/go.uuid"
)

const TodoStream = "todos"

type TodoRepository struct {
	rootRepo eventstore.Repository
}

func (t *TodoRepository) Get(ctx context.Context, todoID uuid.UUID) (*todo.Todo, error) {
	iterator, err := t.rootRepo.GetAggregate(ctx, todoID)
	if err != nil {
		return nil, err
	}
	aggregate := todo.TodoFromHistory(iterator)

	if aggregate.Deleted || aggregate.Version() == 0 {
		return nil, todo.ErrTodoNotFound
	}

	return aggregate, nil
}

func (t *TodoRepository) Save(ctx context.Context, todo *todo.Todo) error {
	return t.rootRepo.SaveAggregate(ctx, todo)
}

func NewTodoRepository(eventStore *eventstore.EventStore) *TodoRepository {
	return &TodoRepository{
		rootRepo: eventstore.NewRepository(TodoStream, eventStore),
	}
}
