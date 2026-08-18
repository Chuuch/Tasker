package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Chuuch/Tasker/backend/internal/apperror"
	"github.com/Chuuch/Tasker/backend/internal/task/model"
)

type mockTaskRepository struct {
	createdTask model.Task
	createCalls int

	updatedTask model.Task
	updateCalls int
	updateErr   error

	getTasks    []model.Task
	getAllCalls int
	getAllErr   error

	getByIDTask  model.Task
	getByIDErr   error
	getByIDCalls int

	deleteTaskCalls int
	deleteTaskError error
}

func (m *mockTaskRepository) Create(
	ctx context.Context,
	task model.Task,
) (model.Task, error) {
	m.createCalls++
	m.createdTask = task
	return task, nil
}

func (m *mockTaskRepository) Update(
	ctx context.Context,
	task model.Task,
) (model.Task, error) {
	m.updateCalls++
	if m.updateErr != nil {
		return model.Task{}, m.updateErr
	}
	m.updatedTask = task
	return task, nil
}

func (m *mockTaskRepository) Delete(ctx context.Context, id int64) error {
	m.deleteTaskCalls++
	return m.deleteTaskError
}

func (m *mockTaskRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	m.getAllCalls++
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.getTasks, nil
}

func (m *mockTaskRepository) GetByID(ctx context.Context, id int64) (model.Task, error) {
	m.getByIDCalls++
	if m.getByIDErr != nil {
		return model.Task{}, m.getByIDErr
	}
	return m.getByIDTask, nil
}

func TestCreateTask(t *testing.T) {
	repository := &mockTaskRepository{}
	service := NewTaskService(repository)

	task, err := service.CreateTask(
		context.Background(),
		"Learn Kubernetes",
		"Understand Pods and Deployments",
	)

	if repository.createCalls != 1 {
		t.Errorf(
			"expected repository Create to be called once, got %d",
			repository.createCalls,
		)
	}

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.Title != "Learn Kubernetes" {
		t.Errorf(
			"expected title %q, got %q",
			"Learn Kubernetes",
			task.Title,
		)
	}

	if task.Description != "Understand Pods and Deployments" {
		t.Errorf(
			"expected description %q, got %q",
			"Understand Pods and Deployments",
			task.Description,
		)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{
			name:  "empty title",
			title: " ",
		},
		{
			name:  "title is too long",
			title: string(make([]byte, 201)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &mockTaskRepository{}
			service := NewTaskService(repository)

			_, err := service.CreateTask(
				context.Background(),
				tt.title,
				"Some description",
			)

			if err == nil {
				t.Fatal("expected an error, got nil")
			}

			if !errors.Is(err, apperror.ErrValidation) {
				t.Errorf(
					"expected error to be apperror.ErrValidation, got %v",
					err,
				)
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	expectedTasks := []model.Task{
		{ID: 1, Title: "Task 1", Description: "Desc 1"},
		{ID: 2, Title: "Task 2", Description: "Desc 2"},
	}

	repository := &mockTaskRepository{
		getTasks: expectedTasks,
	}
	service := NewTaskService(repository)

	tasks, err := service.GetTasks(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.getAllCalls != 1 {
		t.Errorf("expected GetAll to be called once, got %d", repository.getAllCalls)
	}

	if len(tasks) != len(expectedTasks) {
		t.Fatalf("expected %d tasks, got %d", len(expectedTasks), len(tasks))
	}
}

func TestGetTasksError(t *testing.T) {
	repository := &mockTaskRepository{
		getAllErr: errors.New("db connection error"),
	}
	service := NewTaskService(repository)

	_, err := service.GetTasks(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if repository.getAllCalls != 1 {
		t.Errorf("expected GetAll to be called once, got %d", repository.getAllCalls)
	}
}
