package routes

import (
	"taski_backend/internal/handler/handlers"
	"taski_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func SyncRoutes(handler *handlers.SyncHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.JWT)
	r.Post("/events", handler.CreateEvents)
	r.Get("/events", handler.GetEvents)
	return r
}