package todo

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type CommandHandler struct {
	todoRepo TodoCollection
}

func (c *CommandHandler) CreateTodoHandler(ctx context.Context, command CreateTodo) error {
	return c.todoRepo.Save(ctx, NewTodo(command.Title, command.Description, command.Deadline))
}

func (c *CommandHandler) UpdateTodoHandler(ctx context.Context, command UpdateTodo) error {
	todo, err := c.todoRepo.Get(ctx, uuid.FromStringOrNil(command.TodoID))
	if err != nil {
		return err
	}

	todo.Update(command.Title, command.Description, command.Deadline)

	return c.todoRepo.Save(ctx, todo)
}

func (c *CommandHandler) DoTodoHandler(ctx context.Context, command DoTodo) error {
	todo, err := c.todoRepo.Get(ctx, uuid.FromStringOrNil(command.TodoID))
	if err != nil {
		return err
	}

	err = todo.Do()
	if err != nil {
		return err
	}

	return c.todoRepo.Save(ctx, todo)
}

func (c *CommandHandler) UndoTodoHandler(ctx context.Context, command UndoTodo) error {
	todo, err := c.todoRepo.Get(ctx, uuid.FromStringOrNil(command.TodoID))
	if err != nil {
		return err
	}

	err = todo.Undo()
	if err != nil {
		return err
	}

	return c.todoRepo.Save(ctx, todo)
}

func (c *CommandHandler) DeleteHandler(ctx context.Context, command DeleteTodo) error {
	todo, err := c.todoRepo.Get(ctx, uuid.FromStringOrNil(command.TodoID))
	if err != nil {
		return err
	}

	todo.Delete()

	return c.todoRepo.Save(ctx, todo)
}

func NewCommandHandler(repo TodoCollection) *CommandHandler {
	return &CommandHandler{repo}
}
