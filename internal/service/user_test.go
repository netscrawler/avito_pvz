package service_test

import (
	"context"
	"errors"
	"testing"

	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"
	"avito_pvz/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUser_Auth(t *testing.T) {
	t.Parallel()
	user, _ := domain.NewUser("user@example.com", "securePassword123", "moderator")

	tests := []struct {
		name       string
		email      string
		password   string
		want       *string
		wantErr    error
		setupMocks func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator)
	}{
		{
			name:     "auth_success",
			email:    "user@example.com",
			password: "securePassword123",
			want:     stringPtr("token"),
			wantErr:  nil,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "user@example.com").
					Return(user, nil)
				jwt.On("GenerateToken", "user@example.com", "moderator").
					Return("token", nil)
			},
		},
		{
			name:     "user_not_found",
			email:    "nonexistent@example.com",
			password: "password123",
			want:     nil,
			wantErr:  models.ErrUserNotFoud,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "nonexistent@example.com").
					Return(nil, domain.ErrNotFound)
			},
		},
		{
			name:     "incorrect_password",
			email:    "user@example.com",
			password: "wrongPassword",
			want:     nil,
			wantErr:  models.ErrInvalidPassword,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "user@example.com").
					Return(&domain.User{
						Email:        "user@example.com",
						PasswordHash: "hashedPassword123",
						Role:         "user",
					}, nil)
			},
		},
		{
			name:     "internal_error_on_get_user",
			email:    "error@example.com",
			password: "password123",
			want:     nil,
			wantErr:  models.ErrInternal,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "error@example.com").
					Return(nil, errors.New("database error"))
			},
		},
		{
			name:     "token_generation_fails",
			email:    "user@example.com",
			password: "securePassword123",
			want:     nil,
			wantErr:  models.ErrInternalCodeGen,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "user@example.com").
					Return(user, nil)

				jwt.On("GenerateToken", "user@example.com", "moderator").
					Return("", models.ErrInternalCodeGen)
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := service.NewMockUserProvider(t)
			jwt := service.NewMockJWTGenerator(t)

			if tt.setupMocks != nil {
				tt.setupMocks(repo, jwt)
			}

			service := service.NewUserService(repo, jwt)

			got, err := service.Auth(context.Background(), tt.email, tt.password)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, *tt.want, *got)
			}

			repo.AssertExpectations(t)
			jwt.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestUser_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		email      domain.Email
		password   string
		role       domain.Role
		want       *domain.User
		wantErr    error
		setupMocks func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator)
	}{
		{
			name:     "user_created_successfully",
			email:    "user@example.com",
			password: "securePassword123",
			role:     "moderator",
			want: &domain.User{
				Email: "user@example.com",
				Role:  "moderator",
			},
			wantErr: nil,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "user@example.com").
					Return(nil, domain.ErrNotFound)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)

				jwt.On("GenerateToken", "user@example.com", "moderator").
					Return("token", nil).
					Maybe()
			},
		},
		{
			name:     "user_already_exists",
			email:    "existing@example.com",
			password: "password",
			role:     "moderator",
			want:     nil,
			wantErr:  models.ErrUserAlreadyExist,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "existing@example.com").
					Return(&domain.User{Email: "existing@example.com"}, nil)
			},
		},
		{
			name:     "get_by_email_fails_internally",
			email:    "err@example.com",
			password: "123",
			role:     "moderator",
			want:     nil,
			wantErr:  models.ErrInternal,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "err@example.com").
					Return(nil, errors.New("db connection failed"))
			},
		},
		{
			name:     "invalid_email_error",
			email:    "invalid-email",
			password: "pass",
			role:     "moderator",
			want:     nil,
			wantErr:  models.ErrInvalidEmail,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "invalid-email").
					Return(nil, domain.ErrNotFound)
			},
		},
		{
			name:     "create_repo_fails_with_internal_error",
			email:    "newuser@example.com",
			password: "12345",
			role:     "moderator",
			want:     nil,
			wantErr:  models.ErrInternal,
			setupMocks: func(repo *service.MockUserProvider, jwt *service.MockJWTGenerator) {
				repo.On("GetByEmail", mock.Anything, "newuser@example.com").
					Return(nil, domain.ErrNotFound)

				repo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.Email == "newuser@example.com" &&
						u.Role == "moderator" &&
						u.PasswordHash != "" && // хеш должен быть установлен
						u.ID != uuid.Nil // ID должен быть установлен
				})).Return(domain.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		tt := tt // захват переменной
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := service.NewMockUserProvider(t)
			jwt := service.NewMockJWTGenerator(t)

			if tt.setupMocks != nil {
				tt.setupMocks(repo, jwt)
			}

			service := service.NewUserService(repo, jwt)

			got, err := service.Create(context.Background(), tt.email, tt.password, tt.role)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.Email, got.Email)
				require.Equal(t, tt.want.Role, got.Role)
			}

			repo.AssertExpectations(t)
			jwt.AssertExpectations(t)
		})
	}
}
