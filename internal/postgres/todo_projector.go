package postgres

import (
	"github.com/go-event-store/eventstore"
	todo "github.com/go-event-store/example/internal"
	"github.com/go-event-store/pg"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewTodoReadModelProjector(eventStore *eventstore.EventStore, pool *pgxpool.Pool) *eventstore.ReadModelProjector {
	projector := eventstore.NewReadModelProjector("todo_read_model", NewTodoReadModel(pg.NewClient(pool)), eventStore, pg.NewProjectionManager(pool))
	projector.
		FromStream(TodoStream, nil).
		Init(func() interface{} {
			return struct{}{}
		}).
		When(map[string]eventstore.EventHandler{
			"TodoWasCreated": func(state interface{}, event eventstore.DomainEvent) (interface{}, error) {
				projector.ReadModel.Stack("insert", map[string]interface{}{
					"id":          event.AggregateID().String(),
					"title":       event.Payload().(todo.TodoWasCreated).Title,
					"description": event.Payload().(todo.TodoWasCreated).Description,
					"deadline":    event.Payload().(todo.TodoWasCreated).Deadline,
					"updated":     event.CreatedAt(),
					"created":     event.CreatedAt(),
				})

				return state, nil
			},
			"TodoWasUpdated": func(state interface{}, event eventstore.DomainEvent) (interface{}, error) {
				projector.ReadModel.Stack("update", map[string]interface{}{
					"title":       event.Payload().(todo.TodoWasUpdated).Title,
					"description": event.Payload().(todo.TodoWasUpdated).Description,
					"deadline":    event.Payload().(todo.TodoWasUpdated).Deadline,
					"done":        false,
					"updated":     event.CreatedAt(),
				}, map[string]interface{}{
					"id": event.AggregateID().String(),
				})

				return state, nil
			},
			"TodoWasDone": func(state interface{}, event eventstore.DomainEvent) (interface{}, error) {
				projector.ReadModel.Stack("update", map[string]interface{}{
					"done":    true,
					"updated": event.CreatedAt(),
				}, map[string]interface{}{
					"id": event.AggregateID().String(),
				})

				return state, nil
			},
			"TodoWasUndone": func(state interface{}, event eventstore.DomainEvent) (interface{}, error) {
				projector.ReadModel.Stack("update", map[string]interface{}{
					"done":    false,
					"updated": event.CreatedAt(),
				}, map[string]interface{}{
					"id": event.AggregateID().String(),
				})

				return state, nil
			},
			"TodoWasDeleted": func(state interface{}, event eventstore.DomainEvent) (interface{}, error) {
				projector.ReadModel.Stack("remove", map[string]interface{}{
					"id": event.AggregateID().String(),
				})

				return state, nil
			},
		})

	return &projector
}
