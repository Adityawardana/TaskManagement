package usecase

import (
	"context"
	"strings"
	"time"

	domain "taskmanagement/internal/modules/task/domain"
	shareddomain "taskmanagement/pkg/shared/domain"
)

type taskUsecaseImpl struct {
	tasks    TaskRepository
	notifier Notifier
}

// AssignResult represents the result of a task assignment
type AssignResult struct {
	Status     string    `json:"status"`
	TaskID     string    `json:"task_id"`
	AssigneeID string    `json:"assignee_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

// NewTaskUsecase usecase impl constructor
func NewTaskUsecase(tasks TaskRepository, notifier Notifier) TaskUsecase {
	return &taskUsecaseImpl{tasks: tasks, notifier: notifier}
}

func (u *taskUsecaseImpl) CreateTask(ctx context.Context, userID, title, description, status string) (domain.Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return domain.Task{}, shareddomain.NewAppError(400, "INVALID_REQUEST", "title is required", nil)
	}
	status = strings.TrimSpace(status)
	if status == "" {
		status = "todo"
	}
	return u.tasks.CreateTask(ctx, userID, title, description, status)
}

func (u *taskUsecaseImpl) ListTasks(ctx context.Context, userID string, q domain.TaskQuery) ([]domain.Task, error) {
	return u.tasks.ListTasks(ctx, userID, q)
}

func (u *taskUsecaseImpl) GetTask(ctx context.Context, userID, taskID string) (domain.Task, error) {
	return u.tasks.GetTask(ctx, userID, taskID)
}

func (u *taskUsecaseImpl) UpdateTask(ctx context.Context, userID, taskID, title, description, status string) (domain.Task, error) {
	title = strings.TrimSpace(title)
	status = strings.TrimSpace(status)
	if title == "" || status == "" {
		return domain.Task{}, shareddomain.NewAppError(400, "INVALID_REQUEST", "title and status are required", nil)
	}
	return u.tasks.UpdateTask(ctx, userID, taskID, title, description, status)
}

func (u *taskUsecaseImpl) DeleteTask(ctx context.Context, userID, taskID string) error {
	return u.tasks.DeleteTask(ctx, userID, taskID)
}

func (u *taskUsecaseImpl) AssignTask(ctx context.Context, userID, taskID, assigneeID string) (AssignResult, error) {
	err := u.tasks.AssignTask(ctx, userID, taskID, assigneeID, func(ctx context.Context) error {
		return u.notifier.NotifyAssignment(ctx, taskID, assigneeID)
	})
	if err != nil {
		return AssignResult{}, err
	}
	return AssignResult{
		Status:     "assigned",
		TaskID:     taskID,
		AssigneeID: assigneeID,
		AssignedAt: time.Now().UTC(),
	}, nil
}
