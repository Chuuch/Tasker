package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Chuuch/Tasker/backend/internal/apperror"
	"github.com/Chuuch/Tasker/backend/internal/task/model"
	"github.com/Chuuch/Tasker/backend/internal/task/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, title, description string) (model.Task, error)
	GetTasks(ctx context.Context) ([]model.Task, error)
	GetTask(ctx context.Context, id int64) (model.Task, error)
	UpdateTask(ctx context.Context, id int64, title string, description string, completed bool) (model.Task, error)
	DeleteTask(ctx context.Context, id int64) error
}

var ErrTaskNotFound = errors.New("task not found")

type taskService struct {
	repository repository.TaskRepository
}

func NewTaskService(repository repository.TaskRepository) TaskService {
	return &taskService{
		repository: repository,
	}
}

func (s *taskService) GetTasks(ctx context.Context) ([]model.Task, error) {
	return s.repository.GetAll(ctx)
}

func (s *taskService) CreateTask(
	ctx context.Context,
	title string,
	description string,
) (model.Task, error) {

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		return model.Task{}, apperror.ErrValidation
	}

	if len(title) > 200 {
		return model.Task{}, apperror.ErrValidation
	}

	task := model.Task{
		Title:       title,
		Description: description,
	}

	return s.repository.Create(ctx, task)
}

func (s *taskService) UpdateTask(
	ctx context.Context,
	id int64,
	title string,
	description string,
	completed bool,
) (model.Task, error) {

	if id <= 0 {
		return model.Task{}, apperror.ErrValidation
	}

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		return model.Task{}, apperror.ErrValidation
	}

	if description == "" {
		return model.Task{}, apperror.ErrValidation
	}

	if len(title) > 200 {
		return model.Task{}, apperror.ErrValidation
	}

	task := model.Task{
		ID:          id,
		Title:       title,
		Description: description,
		Completed:   completed,
	}

	return s.repository.Update(ctx, task)
}

func (s *taskService) GetTask(ctx context.Context, id int64) (model.Task, error) {
	if id <= 0 {
		return model.Task{}, apperror.ErrValidation
	}

	return s.repository.GetByID(ctx, id)
}

func (s *taskService) DeleteTask(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperror.ErrValidation
	}

	return s.repository.Delete(ctx, id)
}
