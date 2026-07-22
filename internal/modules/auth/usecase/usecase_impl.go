package usecase

import (
	"context"
	"strings"

	domain "taskmanagement/internal/modules/auth/domain"
	shareddomain "taskmanagement/pkg/shared/domain"

	"github.com/google/uuid"
)

type authUsecaseImpl struct {
	users  UserRepository
	hasher PasswordService
	tokens TokenIssuer
}

// NewAuthUsecase usecase impl constructor
func NewAuthUsecase(users UserRepository, hasher PasswordService, tokens TokenIssuer) AuthUsecase {
	return &authUsecaseImpl{users: users, hasher: hasher, tokens: tokens}
}

func (u *authUsecaseImpl) Register(ctx context.Context, email, password, teamName string) (domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || len(password) < 8 {
		return domain.User{}, shareddomain.NewAppError(400, "INVALID_REQUEST", "email and password(min 8 chars) are required", nil)
	}
	hash, err := u.hasher.Hash(password)
	if err != nil {
		return domain.User{}, shareddomain.NewAppError(500, "HASH_ERROR", "failed to process password", err)
	}

	newUserID, err := uuid.NewV7()
	if err != nil {
		return domain.User{}, shareddomain.NewAppError(500, "HASH_ERROR", "failed to generate user id", err)
	}

	userExist, err := u.users.DoesUserIdExist(ctx, newUserID)
	if userExist {
		return domain.User{}, shareddomain.NewAppError(500, "HASH_ERROR", "User Id Already Exist", err)
	}
	if err != nil {
		return domain.User{}, shareddomain.NewAppError(500, "DATABASE_ERROR", "Failed Check User Data", err)
	}

	return u.users.CreateUserWithTeam(ctx, email, hash, teamName, newUserID)
}

func (u *authUsecaseImpl) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", shareddomain.NewAppError(400, "INVALID_REQUEST", "email and password are required", nil)
	}
	user, err := u.users.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if err := u.hasher.Compare(user.PasswordHash, password); err != nil {
		return "", shareddomain.NewAppError(401, "INVALID_CREDENTIALS", "email or password is incorrect", err)
	}
	token, err := u.tokens.Generate(user.ID)
	if err != nil {
		return "", shareddomain.NewAppError(500, "TOKEN_ERROR", "failed to issue token", err)
	}
	return token, nil
}
