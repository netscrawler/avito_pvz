package domain

import (
	"avito_pvz/internal/http/gen"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

type Role string

func (r Role) IsValid() bool {
	return r == RoleModerator || r == RoleEmploye
}

const (
	RoleModerator Role = "moderator"
	RoleEmploye   Role = "employee"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (u *User) CheckPasswordHash(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func NewUser(email string, password string, role string) (*User, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, ErrInvalidEmail
	}
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, ErrInternal
	}
	if role != string(RoleModerator) && role != string(RoleEmploye) {
		return nil, ErrInvalidRole
	}
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         RoleModerator,
		CreatedAt:    time.Now(),
	}, nil
}

func (u *User) ToDto() *gen.User {
	return &gen.User{
		Email: types.Email(u.Email),
		Id:    &u.ID,
		Role:  gen.UserRole(u.Role),
	}
}
