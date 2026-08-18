package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Chuuch/Tasker/backend/internal/apperror"
	"github.com/Chuuch/Tasker/backend/internal/task/service"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetTasks(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get tasks",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	limitRequestBody(w, r)
	var request createTaskRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decodeJSON(r, &request); err != nil {
		writeJSONDecodeError(w, err)
		return
	}

	task, err := h.service.CreateTask(
		r.Context(),
		request.Title,
		request.Description,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	limitRequestBody(w, r)
	idString := r.PathValue("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(
			w,
			"invalid task ID",
			http.StatusBadRequest,
		)
		return
	}

	var request updateTaskRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeJSONDecodeError(w, err)
		return
	}

	task, err := h.service.UpdateTask(
		r.Context(),
		id,
		request.Title,
		request.Description,
		request.Completed,
	)

	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrValidation):
			http.Error(
				w,
				"invalid task data",
				http.StatusBadRequest,
			)

		case errors.Is(err, apperror.ErrNotFound):
			http.Error(
				w,
				"task not found",
				http.StatusNotFound,
			)

		default:
			http.Error(
				w,
				"failed to update task",
				http.StatusInternalServerError,
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(
			w,
			"invalid task ID",
			http.StatusBadRequest,
		)
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrValidation):
			http.Error(
				w,
				"invalid task ID",
				http.StatusBadRequest,
			)

		case errors.Is(err, apperror.ErrNotFound):
			http.Error(
				w,
				"task not found",
				http.StatusNotFound,
			)

		default:
			http.Error(
				w,
				"failed to get task",
				http.StatusInternalServerError,
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(
			w,
			"invalid task ID",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.DeleteTask(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrValidation):
			http.Error(
				w,
				"invalid task ID",
				http.StatusBadRequest,
			)

		case errors.Is(err, apperror.ErrNotFound):
			http.Error(
				w,
				"task not found",
				http.StatusNotFound,
			)

		default:
			http.Error(
				w,
				"failed to delete task",
				http.StatusInternalServerError,
			)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
