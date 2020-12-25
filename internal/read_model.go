package todo

import (
	"context"
	"time"
)

type TodoView struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodoFinder interface {
	Find(ctx context.Context, todoID string) (TodoView, error)
	FindAll(ctx context.Context) ([]TodoView, error)
}
