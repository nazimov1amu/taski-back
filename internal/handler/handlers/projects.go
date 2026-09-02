package handlers

import (
	"encoding/json"
	"net/http"

	"taski_backend/internal/models"
	"taski_backend/internal/service"
)

type ProjectsHandler struct {
	service *service.ProjectsService
}

func NewProjectsHandler(service *service.ProjectsService) *ProjectsHandler {
	return &ProjectsHandler{service: service}
}

func (h *ProjectsHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid_id", http.StatusBadRequest)
		return
	}

	project, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "project")
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectsHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	var args models.ProjectFilters

	query := r.URL.Query()

	name := query.Get("name")
	if name != "" {
		args.Name = &name
	}

	projects, err := h.service.GetBulk(r.Context(), args)
	if err != nil {
		writeServiceError(w, err, "project")
		return
	}

	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectsHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "project")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *ProjectsHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid_id", http.StatusBadRequest)
		return
	}

	var req models.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		writeServiceError(w, err, "project")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *ProjectsHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid_id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err, "project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
