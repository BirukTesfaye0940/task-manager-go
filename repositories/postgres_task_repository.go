package repositories

import (
	"context"
	"database/sql"
	"task-manager-go/models"
)

type PostgresTaskRepo struct {
	db *sql.DB
}

func NewPostgresTaskRepo(db *sql.DB) *PostgresTaskRepo {
	return &PostgresTaskRepo{db: db}
}

func (r *PostgresTaskRepo) GetAll(ctx context.Context, userID int) ([]models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, title, description, done FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepo) GetByID(ctx context.Context, id int, userID int) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}
	var task models.Task
	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, title, description, done FROM tasks WHERE id = $1 AND user_id = $2", id, userID).
		Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}
	return task, nil
}

func (r *PostgresTaskRepo) Create(ctx context.Context, task models.Task) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO tasks (user_id, title, description, done) VALUES ($1, $2, $3, $4) RETURNING id",
		task.UserID, task.Title, task.Description, task.Done,
	).Scan(&task.ID)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (r *PostgresTaskRepo) Update(ctx context.Context, id int, task models.Task) (models.Task, error) {
	if err := ctx.Err(); err != nil {
		return models.Task{}, err
	}
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE tasks SET title = $1, description = $2, done = $3 WHERE id = $4 AND user_id = $5",
		task.Title, task.Description, task.Done, id, task.UserID,
	)
	if err != nil {
		return models.Task{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Task{}, err
	}

	if rowsAffected == 0 {
		return models.Task{}, ErrTaskNotFound
	}

	task.ID = id
	return task, nil
}

func (r *PostgresTaskRepo) Delete(ctx context.Context, id int, userID int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	res, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
