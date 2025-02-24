package models

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CreateTaskRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type UpdateTaskRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
