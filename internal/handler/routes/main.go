package routes

import (
	"taski_backend/internal/handler/handlers"

	"github.com/go-chi/chi/v5"
)

func MainRoutes(tasksHandler *handlers.TasksHandler, projectsHandler *handlers.ProjectsHandler, usersHandler *handlers.UsersHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Mount("/tasks", TasksRoutes(tasksHandler))
		r.Mount("/projects", ProjectsRoutes(projectsHandler))
		r.Mount("/users", UserRoutes(usersHandler))
	})
	return r
}
