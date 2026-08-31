package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"taski_backend/internal/models"
	"taski_backend/internal/service"
)

type TasksHandler struct {
	service *service.TasksService
}

func NewTasksHandler(service *service.TasksService) *TasksHandler {
	return &TasksHandler{service: service}
}

func (h *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var args models.TaskFilters

	query := r.URL.Query()

	dateStr := query.Get("date")
	if dateStr != "" {
		parsed, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			http.Error(w, "invalid date", http.StatusBadRequest)
			return
		}
		dateStr = parsed.Format(time.DateOnly)
		args.Date = &dateStr
	}

	projectID := query.Get("project_id")
	if projectID != "" {
		if _, err := parseUUID(projectID); err != nil {
			http.Error(w, "invalid project_id", http.StatusBadRequest)
			return
		}
		args.ProjectID = &projectID
	}

	title := query.Get("title")
	if title != "" {
		args.Title = &title
	}

	tasks, err := h.service.GetBulk(r.Context(), args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ProjectID != "" {
		if _, err := parseUUID(req.ProjectID); err != nil {
			http.Error(w, "invalid project_id", http.StatusBadRequest)
			return
		}
	}

	created, err := h.service.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ProjectID != "" {
		if _, err := parseUUID(req.ProjectID); err != nil {
			http.Error(w, "invalid project_id", http.StatusBadRequest)
			return
		}
	}

	updated, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
