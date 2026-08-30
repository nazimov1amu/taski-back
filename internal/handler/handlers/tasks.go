package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"taski_backend/internal/handler/dto"
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

	writeJSON(w, http.StatusOK, toTaskResponse(task))
}

func (h *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	date, err := time.Parse(time.DateOnly, r.URL.Query().Get("date"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks, err := h.service.GetBulk(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dto.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		resp = append(resp, toTaskResponse(task))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
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

	created, err := h.service.Create(r.Context(), models.Task{
		Title:       req.Title,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		PlannedAt:   req.PlannedAt,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, toTaskResponse(created))
}

func (h *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := parseUUID(req.ID); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if req.ProjectID != "" {
		if _, err := parseUUID(req.ProjectID); err != nil {
			http.Error(w, "invalid project_id", http.StatusBadRequest)
			return
		}
	}

	updated, err := h.service.Update(r.Context(), models.Task{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		Completed:   req.Completed,
		PlannedAt:   req.PlannedAt,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(updated))
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

func toTaskResponse(task models.Task) dto.TaskResponse {
	return dto.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		ProjectID:   task.ProjectID,
		Completed:   task.Completed,
		PlannedAt:   task.PlannedAt,
		StartTime:   task.StartTime,
		EndTime:     task.EndTime,
	}
}
