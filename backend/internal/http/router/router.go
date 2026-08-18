package router

import (
	"net/http"

	"github.com/Chuuch/Tasker/backend/internal/task/handler"
)

func NewRouter(
	taskHandler *handler.TaskHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /api/v1/tasks", taskHandler.GetAll)
	mux.HandleFunc("GET /api/v1/tasks/{id}", taskHandler.GetByID)
	mux.HandleFunc("POST /api/v1/tasks", taskHandler.Create)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", taskHandler.Update)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", taskHandler.Delete)
	return mux
}
