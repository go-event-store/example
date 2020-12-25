package todo

import "context"

type QueryHandler struct {
	finder TodoFinder
}

func (qh *QueryHandler) FindTodoByIDHandler(ctx context.Context, query FindTodoByID) (TodoView, error) {
	return qh.finder.Find(ctx, query.TodoID)
}

func (qh *QueryHandler) FindAllTodosHandler(ctx context.Context, _ FindAllTodos) ([]TodoView, error) {
	return qh.finder.FindAll(ctx)
}

func NewQueryHandler(finder TodoFinder) *QueryHandler {
	return &QueryHandler{finder}
}
