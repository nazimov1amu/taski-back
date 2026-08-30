package routes

import (
	"taski_backend/internal/handler/handlers"

	"github.com/go-chi/chi/v5"
)

func MainRoutes(tasksHandler *handlers.TasksHandler, projectsHandler *handlers.ProjectsHandler, notesHandler *handlers.NotesHandler) chi.Router {
	r := chi.NewRouter()
	r.Mount("/tasks", TasksRoutes(tasksHandler))
	r.Mount("/projects", ProjectsRoutes(projectsHandler))
	r.Mount("/notes", NotesRoutes(notesHandler))
	return r
}
