package repositories

import (
	"database/sql"
	"task-manager-go/models"
)

type PostgresTaskRepo struct {
	db *sql.DB
}

func NewPostgresTaskRepo(db *sql.DB) *PostgresTaskRepo {
	return &PostgresTaskRepo{db: db}
}

func (r *PostgresTaskRepo) GetAll() ([]models.Task, error) {
	rows, err := r.db.Query("SELECT id, title, description, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepo) GetByID(id int) (models.Task, error) {
	var task models.Task
	err := r.db.QueryRow("SELECT id, title, description, done FROM tasks WHERE id = $1", id).
		Scan(&task.ID, &task.Title, &task.Description, &task.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}
	return task, nil
}

func (r *PostgresTaskRepo) Create(task models.Task) (models.Task, error) {
	err := r.db.QueryRow(
		"INSERT INTO tasks (title, description, done) VALUES ($1, $2, $3) RETURNING id",
		task.Title, task.Description, task.Done,
	).Scan(&task.ID)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (r *PostgresTaskRepo) Update(id int, task models.Task) (models.Task, error) {
	res, err := r.db.Exec(
		"UPDATE tasks SET title = $1, description = $2, done = $3 WHERE id = $4",
		task.Title, task.Description, task.Done, id,
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

func (r *PostgresTaskRepo) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", id)
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
