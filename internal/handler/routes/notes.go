package routes

import (
	"taski_backend/internal/handler/handlers"

	"github.com/go-chi/chi/v5"
)

func NotesRoutes(notesHandler *handlers.NotesHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/bulk", notesHandler.GetNotes)
	r.Get("/", notesHandler.GetNote)
	r.Post("/", notesHandler.CreateNote)
	r.Put("/", notesHandler.UpdateNote)
	r.Delete("/", notesHandler.DeleteNote)
	return r
}
