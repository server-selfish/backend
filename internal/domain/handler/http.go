package handler

import "github.com/go-chi/chi/v5"

func RegisterTodoRoutes(r *chi.Mux, th TodoHandler) {
	r.Post("/todo", th.AddNewTodo)
}
