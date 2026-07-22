package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	domain "taskmanagement/internal/modules/task/domain"
	shareddomain "taskmanagement/pkg/shared/domain"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(ctx context.Context, ownerID, title, description, status string) (domain.Task, error) {
	if strings.TrimSpace(status) == "" {
		status = "todo"
	}
	var t domain.Task
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tasks (owner_id, title, description, status) VALUES ($1, $2, $3, $4)
		 RETURNING id, owner_id, title, description, status, assigned_to, created_at, updated_at`,
		ownerID, title, description, status,
	).Scan(&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.Status, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return domain.Task{}, err
	}
	return t, nil
}

func (r *TaskRepository) ListTasks(ctx context.Context, userID string, q domain.TaskQuery) ([]domain.Task, error) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	offset := (q.Page - 1) * q.Limit

	base := `SELECT id, owner_id, title, description, status, assigned_to, created_at, updated_at
		FROM tasks WHERE (owner_id = $1 OR assigned_to = $1)`
	args := []any{userID}
	argN := 2

	if strings.TrimSpace(q.Status) != "" {
		base += fmt.Sprintf(" AND status = $%d", argN)
		args = append(args, q.Status)
		argN++
	}
	if strings.TrimSpace(q.Search) != "" {
		base += fmt.Sprintf(" AND title ILIKE $%d", argN)
		args = append(args, "%"+q.Search+"%")
		argN++
	}

	base += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argN, argN+1)
	args = append(args, q.Limit, offset)

	rows, err := r.db.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Task, 0)
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.Status, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TaskRepository) GetTask(ctx context.Context, userID, taskID string) (domain.Task, error) {
	var t domain.Task
	err := r.db.QueryRowContext(ctx,
		`SELECT id, owner_id, title, description, status, assigned_to, created_at, updated_at
		 FROM tasks WHERE id = $1 AND (owner_id = $2 OR assigned_to = $2)`,
		taskID, userID,
	).Scan(&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.Status, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, shareddomain.NewAppError(404, "TASK_NOT_FOUND", "task not found", err)
		}
		return domain.Task{}, err
	}
	return t, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, userID, taskID, title, description, status string) (domain.Task, error) {
	var t domain.Task
	err := r.db.QueryRowContext(ctx,
		`UPDATE tasks
		 SET title = $1, description = $2, status = $3, updated_at = NOW()
		 WHERE id = $4 AND (owner_id = $5 OR assigned_to = $5)
		 RETURNING id, owner_id, title, description, status, assigned_to, created_at, updated_at`,
		title, description, status, taskID, userID,
	).Scan(&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.Status, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, shareddomain.NewAppError(404, "TASK_NOT_FOUND", "task not found", err)
		}
		return domain.Task{}, err
	}
	return t, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, ownerID, taskID string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1 AND owner_id = $2`, taskID, ownerID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return shareddomain.NewAppError(404, "TASK_NOT_FOUND", "task not found", sql.ErrNoRows)
	}
	return nil
}

func (r *TaskRepository) AssignTask(ctx context.Context, ownerID, taskID, assigneeID string, notify func(context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentOwner string
	var oldAssigned sql.NullString
	err = tx.QueryRowContext(ctx,
		`SELECT owner_id, assigned_to FROM tasks WHERE id = $1 FOR UPDATE`,
		taskID,
	).Scan(&currentOwner, &oldAssigned)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shareddomain.NewAppError(404, "TASK_NOT_FOUND", "task not found", err)
		}
		return err
	}
	if currentOwner != ownerID {
		return shareddomain.NewAppError(403, "FORBIDDEN", "you are not allowed to assign this task", nil)
	}

	var sameTeam bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM team_members a
			JOIN team_members b ON a.team_id = b.team_id
			WHERE a.user_id = $1 AND b.user_id = $2
		)
	`, ownerID, assigneeID).Scan(&sameTeam)
	if err != nil {
		return err
	}
	if !sameTeam {
		return shareddomain.NewAppError(400, "INVALID_ASSIGNEE", "assignee must be on the same team", nil)
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE tasks SET assigned_to = $1, updated_at = NOW() WHERE id = $2`,
		assigneeID, taskID,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO task_logs (task_id, action, old_value, new_value, performed_by)
		 VALUES ($1, 'assign', $2, $3, $4)`,
		taskID, oldAssigned, assigneeID, ownerID,
	); err != nil {
		return err
	}

	if err := notify(ctx); err != nil {
		return shareddomain.NewAppError(422, "TRANSACTION_FAILED", "failed to assign task", err)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
