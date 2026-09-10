package routes

import (
	"taski_backend/internal/handler/handlers"

	"github.com/go-chi/chi/v5"
)

func MainRoutes(usersHandler *handlers.UsersHandler, syncHandler *handlers.SyncHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Mount("/users", UserRoutes(usersHandler))
		r.Mount("/sync", SyncRoutes(syncHandler))
	})

	return r
}
