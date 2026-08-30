package routes

import (
	"taski_backend/internal/handler/handlers"

	"github.com/go-chi/chi/v5"
)

func TasksRoutes(tasksHandler *handlers.TasksHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/bulk", tasksHandler.GetTasks)
	r.Get("/", tasksHandler.GetTask)
	r.Post("/", tasksHandler.CreateTask)
	r.Put("/", tasksHandler.UpdateTask)
	r.Delete("/", tasksHandler.DeleteTask)
	return r
}
