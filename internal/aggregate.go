package todo

import (
	"context"
	"errors"
	"time"

	"github.com/go-event-store/eventstore"
	uuid "github.com/satori/go.uuid"
)

var ErrTodoNotFound = errors.New("Todo was not found")
var ErrTodoAlreadyDone = errors.New("Todo is already done")
var ErrTodoNotDone = errors.New("Todo is not done")

type TodoCollection interface {
	Get(ctx context.Context, todoID uuid.UUID) (*Todo, error)
	Save(ctx context.Context, todo *Todo) error
}

type TodoWasCreated struct {
	Title       string
	Description string
	Deadline    time.Time
}

type TodoWasUpdated struct {
	Title       string
	Description string
	Deadline    time.Time
}

type TodoWasDone struct{}
type TodoWasUndone struct{}
type TodoWasDeleted struct{}

type Todo struct {
	eventstore.BaseAggregate

	Title       string
	Description string
	Done        bool
	Deleted     bool
	Deadline    time.Time
}

func (t *Todo) TodoID() uuid.UUID {
	return t.AggregateID()
}

func (t *Todo) Do() error {
	if t.Done {
		return ErrTodoAlreadyDone
	}

	t.RecordThat(TodoWasDone{}, nil)
	return nil
}

func (t *Todo) Undo() error {
	if t.Done == false {
		return ErrTodoNotDone
	}

	t.RecordThat(TodoWasUndone{}, nil)
	return nil
}

func (t *Todo) Update(title, description string, deadline time.Time) {
	t.RecordThat(TodoWasUpdated{title, description, deadline}, map[string]interface{}{})
}

func (t *Todo) Delete() {
	t.RecordThat(TodoWasDeleted{}, nil)
}

func (t *Todo) WhenTodoWasCreated(e TodoWasCreated, _ map[string]interface{}) {
	t.Title = e.Title
	t.Description = e.Description
	t.Deadline = e.Deadline
	t.Done = false
	t.Deleted = false
}

func (t *Todo) WhenTodoWasUpdated(e TodoWasUpdated, _ map[string]interface{}) {
	t.Title = e.Title
	t.Description = e.Description
	t.Deadline = e.Deadline
	t.Done = false
}

func (t *Todo) WhenTodoWasDone(_ TodoWasDone, _ map[string]interface{}) {
	t.Done = true
}

func (t *Todo) WhenTodoWasUndone(_ TodoWasUndone, _ map[string]interface{}) {
	t.Done = false
}

func (t *Todo) WhenTodoWasDeleted(_ TodoWasDeleted, _ map[string]interface{}) {
	t.Deleted = true
}

func NewTodo(title, description string, deadline time.Time) *Todo {
	todo := new(Todo)
	todo.BaseAggregate = eventstore.NewAggregate(todo)
	todo.RecordThat(TodoWasCreated{title, description, deadline}, map[string]interface{}{})

	return todo
}

func TodoFromHistory(events eventstore.DomainEventIterator) *Todo {
	todo := new(Todo)
	todo.BaseAggregate = eventstore.NewAggregate(todo)
	todo.FromHistory(events)

	return todo
}
