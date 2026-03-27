package handler

import (
	"net/http"

	"github.com/server-selfish/backend/internal/domain/service"
)

type (
	TodoHandler interface {
		// your function definition
		AddNewTodo(w http.ResponseWriter, r *http.Request)
	}
	todoHandler struct {
		// your injected dependency
		// for example
		ts service.TodoService
	}
)

func NewTodoHandler(ts service.TodoService) TodoHandler {
	return &todoHandler{
		ts: ts,
	}
}

func (t *todoHandler) AddNewTodo(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
