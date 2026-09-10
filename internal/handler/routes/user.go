package routes

import (
	"taski_backend/internal/handler/handlers"
	"taski_backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)


func UserRoutes(handler *handlers.UsersHandler) chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT)
		r.Get("/me", handler.GetUser)
		r.Put("/", handler.UpdateUser)
		r.Delete("/", handler.DeleteUser)
	})
	r.Group(func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.RefreshToken)
		r.Post("/code", handler.VerifyCode)
	})
	return r
}