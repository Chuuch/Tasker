package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Chuuch/Tasker/backend/internal/apperror"
	"github.com/Chuuch/Tasker/backend/internal/task/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	databaseURL := "postgres://tasker:tasker_password@localhost:5432/tasker_db_test?sslmode=disable"

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("failed to ping database: %v", err)
	}

	_, err = pool.Exec(
		ctx,
		"TRUNCATE TABLE tasks RESTART IDENTITY",
	)
	if err != nil {
		pool.Close()
		t.Fatalf("failed to clean test database: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestPostgresTaskRepository_Create(t *testing.T) {
	pool := setupTestDatabase(t)

	repository := NewPostgresTaskRepository(pool)
	ctx := context.Background()

	task := model.Task{
		Title:       "Integration Test Task",
		Description: "Created using a real PostgreSQL database",
	}

	createdTask, err := repository.Create(ctx, task)

	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	if createdTask.ID == 0 {
		t.Errorf("expected task ID to be generated")
	}

	if createdTask.Title != task.Title {
		t.Errorf(
			"expected title %q, got %q",
			task.Title,
			createdTask.Title,
		)
	}

	if createdTask.Description != task.Description {
		t.Errorf(
			"expected descriptoin %q, got %q",
			task.Description,
			createdTask.Description,
		)
	}

	if createdTask.Completed {
		t.Error("expected new task to be incomplete")
	}

	if createdTask.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be populated")
	}

	if createdTask.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be populated")
	}
}

func TestPostgresTaskRepository_Update(t *testing.T) {
	pool := setupTestDatabase(t)

	repository := NewPostgresTaskRepository(pool)

	ctx := context.Background()

	createdTask, err := repository.Create(ctx, model.Task{
		Title:       "Original Title",
		Description: "Original Description",
	})

	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	createdAt := createdTask.CreatedAt

	updatedTask, err := repository.Update(ctx, model.Task{
		ID:          createdTask.ID,
		Title:       "Updated Title",
		Description: "Updated Description",
		Completed:   true,
	})

	if err != nil {
		t.Fatalf("failed to udpated task: %v", err)
	}

	if updatedTask.ID != createdTask.ID {
		t.Errorf(
			"expected ID %d, got %d",
			createdTask.ID,
			updatedTask.ID,
		)
	}

	if updatedTask.Title != "Updated Title" {
		t.Errorf(
			"expected title %q, got %q",
			"Updated Task",
			updatedTask.Title,
		)
	}

	if updatedTask.Description != "Updated Description" {
		t.Error(
			"expected description",
			"Updated Description",
			updatedTask.Description,
		)
	}

	if !updatedTask.Completed {
		t.Error("expected task to be completed")
	}

	if !updatedTask.CreatedAt.Equal(createdAt) {
		t.Error("expected CreatedAt to remain unchanged")
	}

	if updatedTask.UpdatedAt.Equal(updatedTask.CreatedAt) {
		t.Error("expected UpdatedAt to be after CreatedAt")
	}
}

func TestPostgresTaskRepository_Delete(t *testing.T) {
	pool := setupTestDatabase(t)
	repository := NewPostgresTaskRepository(pool)
	ctx := context.Background()

	task, err := repository.Create(ctx, model.Task{
		Title:       "Task to delete",
		Description: "Delete me",
	})

	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	err = repository.Delete(ctx, task.ID)

	if err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}

	_, err = repository.GetByID(ctx, task.ID)

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf(
			"expected ErrNotFound, got %v",
			err,
		)
	}
}

func TestPostgresTaskRepository_DeleteNotFound(t *testing.T) {
	pool := setupTestDatabase(t)
	repository := NewPostgresTaskRepository(pool)

	err := repository.Delete(
		context.Background(),
		99999,
	)

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf(
			"expected ErrNotFound, got %v",
			err,
		)
	}
}
