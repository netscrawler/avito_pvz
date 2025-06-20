package service

import (
	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"
	"context"
	"errors"
)

type JWTGenerator interface {
	GenerateToken(email, role string) (string, error)
}
type UserProvider interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type User struct {
	jwt  JWTGenerator
	repo UserProvider
}

func (u *User) Create(
	ctx context.Context,
	email domain.Email,
	password string,
	role domain.Role,
) (*domain.User, error) {
	_, err := u.repo.GetByEmail(ctx, string(email))
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrInternal
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrUserAlreadyExist
	}

	user, err := domain.NewUser(string(email), password, string(role))
	if errors.Is(err, domain.ErrInvalidEmail) {
		return nil, models.ErrInvalidEmail
	}

	if err != nil {
		return nil, models.ErrInternal
	}

	err = u.repo.Create(ctx, user)
	if err != nil {
		return nil, models.ErrInternal
	}

	return user, nil
}

func (u *User) Auth(ctx context.Context, email string, password string) (*string, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrUserNotFoud
	}

	if err != nil {
		return nil, models.ErrInternal
	}

	if !user.CheckPasswordHash(password) {
		return nil, models.ErrInvalidPassword
	}

	token, err := u.jwt.GenerateToken(user.Email, string(user.Role))
	if err != nil {
		return nil, models.ErrInternalCodeGen
	}

	return &token, nil
}

func NewUserService(repo UserProvider, jwt JWTGenerator) *User {
	return &User{
		repo: repo,
		jwt:  jwt,
	}
}
