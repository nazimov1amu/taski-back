package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"taski_backend/internal/config"
	"taski_backend/internal/db"
	"taski_backend/internal/handler/handlers"
	"taski_backend/internal/handler/routes"
	"taski_backend/internal/repository"
	"taski_backend/internal/service"
)

type App struct {
	config *config.AppConfig
	db     *sql.DB
	router chi.Router
}

func New() (*App, error) {
	cfg, err := config.NewAppConfig()
	if err != nil {
		return nil, err
	}
	if cfg.Address == "" {
		cfg.Address = ":8080"
	}

	sqlDB, err := db.NewDB(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	tasksHandler := handlers.NewTasksHandler(service.NewTasksService(repository.NewTasksRepository(sqlDB)))
	projectsHandler := handlers.NewProjectsHandler(service.NewProjectsService(repository.NewProjectsRepository(sqlDB)))
	notesHandler := handlers.NewNotesHandler(service.NewNotesService(repository.NewNotesRepository(sqlDB)))

	router := routes.MainRoutes(tasksHandler, projectsHandler, notesHandler)

	return NewApp(cfg, sqlDB, router), nil
}

func NewApp(cfg *config.AppConfig, db *sql.DB, router chi.Router) *App {
	return &App{config: cfg, db: db, router: router}
}

func (a *App) Addr() string {
	return a.config.Address
}

func (a *App) Run() error {
	return http.ListenAndServe(a.config.Address, a.router)
}

func (a *App) Close() error {
	if a.db == nil {
		return nil
	}
	return a.db.Close()
}
