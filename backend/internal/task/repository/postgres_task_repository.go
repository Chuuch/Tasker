package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Chuuch/Tasker/backend/internal/apperror"
	"github.com/Chuuch/Tasker/backend/internal/task/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTaskRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTaskRepository(db *pgxpool.Pool) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: db,
	}
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task model.Task) (model.Task, error) {
	query := `
			UPDATE tasks
			SET
					title = $1,
					description = $2,
					completed = $3,
					updated_at = NOW()
			WHERE id = $4
			RETURNING
					id,
					title,
					description,
					completed,
					created_at,
					updated_at
	`

	var updatedTask model.Task

	err := r.db.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.ID,
	).Scan(
		&updatedTask.ID,
		&updatedTask.Title,
		&updatedTask.Description,
		&updatedTask.Completed,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, apperror.ErrNotFound
		}

		return model.Task{}, fmt.Errorf("updating task %d: %w", task.ID, err)
	}

	return updatedTask, nil
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task model.Task) (model.Task, error) {
	query := `
		INSERT INTO tasks (
			title,
			description
	)
	VALUES ($1, $2)
	RETURNING
		id,
		title,
		description,
		completed,
		created_at,
		updated_at
	`

	err := r.db.QueryRow(ctx, query, task.Title, task.Description).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return model.Task{}, fmt.Errorf("creating task: %w", err)
	}

	return task, nil
}

func (r *PostgresTaskRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	query := `
			SELECT
					id,
					title,
					description,
					completed,
					created_at,
					updated_at
			FROM tasks
			ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("getting tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning task: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tasks: %w", err)
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id int64) (model.Task, error) {
	query := `
			SELECT
					id,
					title,
					description,
					completed,
					created_at,
					updated_at
			FROM tasks
			WHERE id = $1
	`

	var task model.Task

	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, apperror.ErrNotFound
		}

		return model.Task{}, fmt.Errorf("getting task %d: %w", id, err)
	}

	return task, nil
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM tasks
		WHERE ID = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting task %d: %w", id, err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
