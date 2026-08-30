package routes

import (
	"taski_backend/internal/handler/handlers"

	"github.com/go-chi/chi/v5"
)

func ProjectsRoutes(projectsHandler *handlers.ProjectsHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/bulk", projectsHandler.GetProjects)
	r.Get("/", projectsHandler.GetProject)
	r.Post("/", projectsHandler.CreateProject)
	r.Put("/", projectsHandler.UpdateProject)
	r.Delete("/", projectsHandler.DeleteProject)
	return r
}
