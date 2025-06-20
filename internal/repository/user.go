package repository

import (
	"avito_pvz/internal/models/domain"
	"context"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type User struct {
	UserRepository
}

func NewUser(u UserRepository) *User {
	return &User{
		UserRepository: u,
	}
}
