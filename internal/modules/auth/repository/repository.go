package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	domain "taskmanagement/internal/modules/auth/domain"
	shareddomain "taskmanagement/pkg/shared/domain"

	"github.com/google/uuid"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUserWithTeam(ctx context.Context, email, passwordHash, teamName string, userId uuid.UUID) (domain.User, error) {
	teamName = strings.TrimSpace(teamName)
	if teamName == "" {
		teamName = "default"
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback()

	var user domain.User
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING id, email, password_hash, created_at`,
		userId, email, passwordHash,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return domain.User{}, shareddomain.NewAppError(409, "EMAIL_EXISTS", "email already registered", err)
		}
		return domain.User{}, err
	}

	var teamID string
	err = tx.QueryRowContext(ctx, `SELECT id FROM teams WHERE name = $1`, teamName).Scan(&teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = tx.QueryRowContext(ctx, `INSERT INTO teams (name) VALUES ($1) RETURNING id`, teamName).Scan(&teamID)
			if err != nil {
				return domain.User{}, err
			}
		} else {
			return domain.User{}, err
		}
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO team_members (team_id, user_id) VALUES ($1, $2) ON CONFLICT (team_id, user_id) DO NOTHING`,
		teamID, user.ID,
	); err != nil {
		return domain.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, shareddomain.NewAppError(401, "INVALID_CREDENTIALS", "email or password is incorrect", err)
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *AuthRepository) DoesUserIdExist(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT id FROM users WHERE id = $1;`

	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
