package todo

import "time"

type CreateTodo struct {
	Title       string    `binding:"required"`
	Description string    `binding:"required"`
	Deadline    time.Time `binding:"required"`
}

type UpdateTodo struct {
	TodoID      string    `json:"todo_id" binding:"required,uuid4"`
	Title       string    `binding:"required"`
	Description string    `binding:"required"`
	Deadline    time.Time `binding:"required"`
}

type DoTodo struct {
	TodoID string `json:"todo_id" binding:"required,uuid4"`
}

type UndoTodo struct {
	TodoID string `json:"todo_id" binding:"required,uuid4"`
}

type DeleteTodo struct {
	TodoID string `json:"todo_id" binding:"required,uuid4"`
}
