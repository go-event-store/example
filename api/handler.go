package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	todo "github.com/go-event-store/example/internal"
)

// TodoHandler API Actions
type TodoHandler struct {
	ch *todo.CommandHandler
	qh *todo.QueryHandler
}

// GetHandler API
// @Summary Get Todo
// @Description Get Single Todo
// @Tags Read Todo
// @Accept  json
// @Produce json
// @Router /todo/{id} [get]
// @Param id path string true "TodoID"
// @Success 200 {object} todo.TodoView
// @Success 404
func (h *TodoHandler) GetHandler(c *gin.Context) {
	view, err := h.qh.FindTodoByIDHandler(c, todo.FindTodoByID{TodoID: c.Param("id")})
	if err == todo.ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, view)
}

// ListHandler API
// @Summary List Todos
// @Description List all Todo
// @Tags Read Todo
// @Accept  json
// @Produce json
// @Router /todo [get]
// @Success 200 {array} todo.TodoView
// @Success 404
func (h *TodoHandler) ListHandler(c *gin.Context) {
	view, err := h.qh.FindAllTodosHandler(c, todo.FindAllTodos{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, view)
}

// CreateHandler API
// @Summary Create Todo
// @Description Create a new Todo
// @Tags Write Todo
// @Accept  json
// @Produce json
// @Router /create-todo [post]
// @Param command body todo.CreateTodo true "create todo"
// @Success 204
// @Success 404
// @Success 500
func (h *TodoHandler) CreateHandler(c *gin.Context) {
	var command todo.CreateTodo

	err := c.ShouldBindBodyWith(&command, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = h.ch.CreateTodoHandler(c, command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UpdateHandler API
// @Summary Update Todo
// @Description Update Todo
// @Tags Write Todo
// @Accept  json
// @Produce json
// @Router /update-todo [post]
// @Param command body todo.UpdateTodo true "update todo"
// @Success 204
// @Success 404
// @Success 500
func (h *TodoHandler) UpdateHandler(c *gin.Context) {
	var command todo.UpdateTodo

	err := c.ShouldBindBodyWith(&command, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = h.ch.UpdateTodoHandler(c, command)
	if err == todo.ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DoHandler API
// @Summary Do Todo
// @Description Do Todo
// @Tags Write Todo
// @Accept  json
// @Produce json
// @Router /do-todo [post]
// @Param command body todo.DoTodo true "do todo"
// @Success 204
// @Success 404
// @Success 500
func (h *TodoHandler) DoHandler(c *gin.Context) {
	var command todo.DoTodo

	err := c.ShouldBindBodyWith(&command, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = h.ch.DoTodoHandler(c, command)
	if err == todo.ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UndoHandler API
// @Summary Undo Todo
// @Description Undo Todo
// @Tags Write Todo
// @Accept  json
// @Produce json
// @Router /undo-todo [post]
// @Param command body todo.UndoTodo true "undo todo"
// @Success 204
// @Success 404
// @Success 500
func (h *TodoHandler) UndoHandler(c *gin.Context) {
	var command todo.UndoTodo

	err := c.ShouldBindBodyWith(&command, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = h.ch.UndoTodoHandler(c, command)
	if err == todo.ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeleteHandler API
// @Summary Delete Todo
// @Description Delete Todo
// @Tags Write Todo
// @Accept  json
// @Produce json
// @Router /delete-todo [post]
// @Param command body todo.UndoTodo true "delete todo"
// @Success 204
// @Success 404
// @Success 500
func (h *TodoHandler) DeleteHandler(c *gin.Context) {
	var command todo.DeleteTodo

	err := c.ShouldBindBodyWith(&command, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = h.ch.DeleteHandler(c, command)
	if err == todo.ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(ch *todo.CommandHandler, qh *todo.QueryHandler) TodoHandler {
	return TodoHandler{ch, qh}
}
