package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chuuch/Tasker/backend/internal/apperror"
	"github.com/Chuuch/Tasker/backend/internal/task/model"
)

type mockTaskService struct {
	createdTitle       string
	createdDescription string
	createCalls        int
	createdTask model.Task

	getTasksCalls int
	tasks         []model.Task

	getTaskCalls  int
	task          model.Task
	getTaskResult model.Task

	updateTaskCalls int
	updatedTask     model.Task

	deleteTaskCalls int

	serviceError error
}

func (m *mockTaskService) CreateTask(
	ctx context.Context,
	title string,
	description string,
) (model.Task, error) {
	m.createCalls++
	m.createdTitle = title
	m.createdDescription = description

	return model.Task{
		ID:          1,
		Title:       title,
		Description: description,
		Completed:   false,
	}, nil
}

func (m *mockTaskService) UpdateTask(
	ctx context.Context,
	id int64,
	title string,
	description string,
	completed bool,
) (model.Task, error) {
	m.updateTaskCalls++
	if m.serviceError != nil {
		return model.Task{}, m.serviceError
	}
	m.updatedTask = model.Task{
		ID:          id,
		Title:       title,
		Description: description,
		Completed:   completed,
	}
	return m.updatedTask, nil
}

func (m *mockTaskService) DeleteTask(ctx context.Context, id int64) error {
	m.deleteTaskCalls++
	return m.serviceError
}

func (m *mockTaskService) GetTasks(ctx context.Context) ([]model.Task, error) {
	m.getTasksCalls++
	return m.tasks, nil
}

func (m *mockTaskService) GetTask(ctx context.Context, id int64) (model.Task, error) {
	m.getTaskCalls++
	if m.serviceError != nil {
		return model.Task{}, m.serviceError
	}
	return m.getTaskResult, nil
}

func TestCreateTask(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	requestBody := `{
		"title": "Learn Kubernetes",
		"description": "Understand Pods and Deployments"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/tasks",
		bytes.NewBufferString(requestBody),
	)

	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.Create(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			response.Code,
		)
	}

	var task model.Task

	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&task); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.Title != "Learn Kubernetes" {
		t.Errorf(
			"expected title %q, got %q",
			"Learn Kubernetes",
			task.Title,
		)
	}

	if service.createCalls != 1 {
		t.Fatalf(
			"expected service to be called once, got %d",
			service.createCalls,
		)
	}

	if service.createdTitle != "Learn Kubernetes" {
		t.Errorf(
			"expected title %q, got %q",
			"Learn Kubernetes",
			service.createdTitle,
		)
	}
}

func TestCreateTaskInvalidJSON(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/tasks",
		bytes.NewBufferString(`{"title":`),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if service.createCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			service.createCalls,
		)
	}
}

func TestGetAllTasks(t *testing.T) {
	service := &mockTaskService{
		tasks: []model.Task{
			{
				ID:          1,
				Title:       "Learn Go",
				Description: "Learn Go testing",
				Completed:   false,
			},
			{
				ID:          2,
				Title:       "Learn Kubernetes",
				Description: "Learn Pods",
				Completed:   false,
			},
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks",
		nil,
	)

	response := httptest.NewRecorder()

	handler.GetAll(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if service.getTasksCalls != 1 {
		t.Fatalf(
			"expected service to be called once, got %d",
			service.getTasksCalls,
		)
	}

	var tasks []model.Task

	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&tasks); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf(
			"expected 2 tasks, got %d",
			len(tasks),
		)
	}

	if tasks[0].Title != "Learn Go" {
		t.Errorf(
			"expected first task %q, got %q",
			"Learn Go",
			tasks[0].Title,
		)
	}
}

func TastGetTask(t *testing.T) {
	service := &mockTaskService{
		task: model.Task{
			ID:          42,
			Title:       "Learn GO",
			Description: "Learn Testing",
			Completed:   false,
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/42",
		nil,
	)

	request.SetPathValue("id", "42")
	response := httptest.NewRecorder()
	handler.GetByID(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if service.getTaskCalls != 1 {
		t.Fatalf(
			"expected service to be called once, got %d",
			service.getTaskCalls,
		)
	}

	var task model.Task

	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&task); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.ID != 42 {
		t.Errorf("expected ID 42, got %d", task.ID)
	}
}

func TestGetTaskInvalidID(t *testing.T) {
	service := &mockTaskService{}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/abc",
		nil,
	)

	request.SetPathValue("id", "abc")
	response := httptest.NewRecorder()
	handler.GetByID(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if service.getTaskCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			service.getTaskCalls,
		)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	service := &mockTaskService{
		serviceError: apperror.ErrNotFound,
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/999",
		nil,
	)

	request.SetPathValue("id", "999")
	response := httptest.NewRecorder()
	handler.GetByID(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.Code,
		)
	}
}

func TestUpdateTask(t *testing.T) {
	service := &mockTaskService{}

	handler := NewTaskHandler(service)

	requestBody := `{
		"title": "Updated Task",
		"description": "Updated Description",
		"completed": true
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/42",
		bytes.NewBufferString(requestBody),
	)

	request.Header.Set("Content-Type", "application/json")
	request.SetPathValue("id", "42")

	response := httptest.NewRecorder()

	handler.Update(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if service.updateTaskCalls != 1 {
		t.Fatalf(
			"expected service to be called once, got %d",
			service.updateTaskCalls,
		)
	}

	if service.updatedTask.ID != 42 {
		t.Fatalf(
			"expected ID 42, got %d",
			service.updatedTask.ID,
		)
	}

	if service.updatedTask.Title != "Updated Task" {
		t.Errorf(
			"expected title %q, got %q",
			"Updated Task",
			service.updatedTask.Title,
		)
	}

	if !service.updatedTask.Completed {
		t.Error("expected task to be completed")
	}
}

func TestUpdateTaskInvalidJSON(t *testing.T) {
	service := &mockTaskService{}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/42",
		bytes.NewBufferString(`{"title":`),
	)

	request.SetPathValue("id", "42")
	response := httptest.NewRecorder()

	handler.Update(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if service.updateTaskCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			response.Code,
		)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	service := &mockTaskService{
		serviceError: apperror.ErrNotFound,
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/999",
		bytes.NewBufferString(`{
			"title": "Test",
			"description": "Test",
			"completed": false
			}`),
	)

	request.SetPathValue("id", "999")
	response := httptest.NewRecorder()

	handler.Update(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.Code,
		)
	}
}

func TestDeleteTask(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/tasks/42",
		nil,
	)

	request.SetPathValue("id", "42")
	response := httptest.NewRecorder()
	handler.Delete(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			response.Code,
		)
	}

	if service.deleteTaskCalls != 1 {
		t.Fatalf(
			"expected service to be called once, got %d",
			service.deleteTaskCalls,
		)
	}

	if response.Body.Len() != 0 {
		t.Errorf(
			"expected empty response body, got %q",
			response.Body.String(),
		)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	service := &mockTaskService{
		serviceError: apperror.ErrNotFound,
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/tasks/9999",
		nil,
	)

	request.SetPathValue("id", "9999")
	response := httptest.NewRecorder()

	handler.Delete(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			response.Code,
		)
	}
}

func TestDeleteTaskInvalidID(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/tasks/abc",
		nil,
	)

	request.SetPathValue("id", "abc")
	response := httptest.NewRecorder()

	handler.Delete(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if service.deleteTaskCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			service.deleteTaskCalls,
		)
	}
}

func TestCreateTaskValidJSON(t *testing.T) {
	service := &mockTaskService{
		createdTask: model.Task{
			ID: 1,
			Title: "Learn Go",
			Description: "Learn Kubernetes",
			Completed: false,
		},
	}

	handler := NewTaskHandler(service)

	body := strings.NewReader(`{
		"title": "Learn Go",
		"description": "Learn Kubernetes"
		}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		body,
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	handler.Create(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			response.Code,
		)
	}

	if service.createCalls != 1 {
		t.Fatalf(
			"expected service to be called once, got %d",
			service.createCalls,
		)
	}
}


func TestCreateTaskMalformedJSON(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	body := strings.NewReader(`{
		"title": "Learn Go",
		"description":
	}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		body,
	)

	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if service.createCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			service.createCalls,
		)
	}
}

func TestCreateTaskUnknownField(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	body := strings.NewReader(`{
		"title": "Learn Go",
		"description": "Kubernetes",
		"unknownField": "this should fail"
	}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		body,
	)
	
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	handler.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if service.createCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			service.createCalls,
		)
	}
}

func TestCreateTAskMultipleJSONValues(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	body := strings.NewReader(
		`{"title": "Learn Go"}{"title": "Learn Kubernetes"}`,
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		body,
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if service.createCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			service.createCalls,
		)
	}
}

func TestCreateTaskOversizedBody(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	largeDescription := strings.Repeat(
		"a",
		maxRequestBodySize+1,
	)

	body := strings.NewReader(
		`{"title": "Large Task","description":"}` +
		largeDescription +
		`"}"`,
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		body,
	)

	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.Create(response, request)

	if response.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusRequestEntityTooLarge,
			response.Code,
		)
	}

	if service.createCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			service.createCalls,
		)
	}
}

func TestUpdateTaskUnknownField(t *testing.T) {
	service := &mockTaskService{}
	handler := NewTaskHandler(service)

	body := strings.NewReader(`{
		"title": "Learn Go",
		"description": "Kubernetes",
		"unknownField": true
	}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		body,
	)

	request.SetPathValue("id", "4")
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler.Update(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if service.updateTaskCalls != 0 {
		t.Fatalf(
			"expected service not to be called, got %d calls",
			service.updateTaskCalls,
		)
	}
}
