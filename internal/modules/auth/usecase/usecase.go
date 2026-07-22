package usecase

import (
	"context"

	domain "taskmanagement/internal/modules/auth/domain"

	"github.com/google/uuid"
)

// AuthUsecase abstraction
type AuthUsecase interface {
	Register(ctx context.Context, email, password, teamName string) (domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type UserRepository interface {
	CreateUserWithTeam(ctx context.Context, email, passwordHash, teamName string, userId uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	DoesUserIdExist(ctx context.Context, userId uuid.UUID) (bool, error)
}

type TokenIssuer interface {
	Generate(userID string) (string, error)
}

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
