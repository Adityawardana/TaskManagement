package domain

import "time"

type Task struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssignedTo  *string   `json:"assigned_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskQuery struct {
	Status string
	Search string
	Limit  int
	Page   int
}

type CreateTaskRequestObject struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateTaskRequestObject struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type AssignTaskRequestObject struct {
	AssigneeID string `json:"assignee_id"`
}

// AssignTaskResponseObject represents the result of a task assignment
type AssignTaskResponseObject struct {
	Status     string    `json:"status"`
	TaskID     string    `json:"task_id"`
	AssigneeID string    `json:"assignee_id"`
	AssignedAt time.Time `json:"assigned_at"`
}
