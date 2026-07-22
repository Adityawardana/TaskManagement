package helper

import (
	"context"
	"log"
)

type Notifier interface {
	NotifyAssignment(ctx context.Context, taskID, assigneeID string) error
}

type MockNotifier struct{}

func NewMockNotifier() *MockNotifier {
	return &MockNotifier{}
}

func (n *MockNotifier) NotifyAssignment(ctx context.Context, taskID, assigneeID string) error {
	log.Printf("notification: task %s assigned to user %s", taskID, assigneeID)
	return nil
}
