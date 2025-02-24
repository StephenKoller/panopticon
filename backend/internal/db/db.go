package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/stephen/panopticon/internal/models"
)

type Database struct {
	*sql.DB
}

func NewDatabase() (*Database, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) InitializeTables() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating tasks table: %v", err)
		return err
	}

	return nil
}

func (db *Database) CreateTask(task *models.CreateTaskRequest) (*models.Task, error) {
	query := `
		INSERT INTO tasks (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed
	`

	newTask := &models.Task{}
	err := db.QueryRow(query, task.Title, task.Completed).Scan(
		&newTask.ID,
		&newTask.Title,
		&newTask.Completed,
	)

	if err != nil {
		return nil, err
	}

	return newTask, nil
}

func (db *Database) GetTasks() ([]models.Task, error) {
	query := `
		SELECT id, title, completed
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (db *Database) UpdateTask(id int, task *models.UpdateTaskRequest) error {
	query := `
		UPDATE tasks
		SET title = $1, completed = $2
		WHERE id = $3
	`

	result, err := db.Exec(query, task.Title, task.Completed, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

func (db *Database) DeleteTask(id int) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}
