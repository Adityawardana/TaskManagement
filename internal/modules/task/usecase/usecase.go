package usecase

import (
	"context"

	domain "taskmanagement/internal/modules/task/domain"
)

// TaskUsecase abstraction
type TaskUsecase interface {
	CreateTask(ctx context.Context, userID, title, description, status string) (domain.Task, error)
	ListTasks(ctx context.Context, userID string, q domain.TaskQuery) ([]domain.Task, error)
	GetTask(ctx context.Context, userID, taskID string) (domain.Task, error)
	UpdateTask(ctx context.Context, userID, taskID, title, description, status string) (domain.Task, error)
	DeleteTask(ctx context.Context, userID, taskID string) error
	AssignTask(ctx context.Context, userID, taskID, assigneeID string) (AssignResult, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, ownerID, title, description, status string) (domain.Task, error)
	ListTasks(ctx context.Context, userID string, q domain.TaskQuery) ([]domain.Task, error)
	GetTask(ctx context.Context, userID, taskID string) (domain.Task, error)
	UpdateTask(ctx context.Context, userID, taskID, title, description, status string) (domain.Task, error)
	DeleteTask(ctx context.Context, ownerID, taskID string) error
	AssignTask(ctx context.Context, ownerID, taskID, assigneeID string, notify func(context.Context) error) error
}

type Notifier interface {
	NotifyAssignment(ctx context.Context, taskID, assigneeID string) error
}
