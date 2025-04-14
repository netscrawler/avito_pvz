package domain_test

import (
	"testing"
	"time"

	"avito_pvz/internal/models/domain"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string // описание теста
		email    string
		password string
		role     string
		want     *domain.User
		wantErr  bool
	}{
		{
			name:     "valid_user_creation",
			email:    "user@example.com",
			password: "securePassword123",
			role:     "moderator",
			want: &domain.User{
				Email:        "user@example.com",
				PasswordHash: "hashedPassword123",
				Role:         "moderator",
				CreatedAt: time.Now().
					Truncate(time.Second),
			},
			wantErr: false,
		},
		{
			name:     "invalid_email_format",
			email:    "invalid-email",
			password: "securePassword123",
			role:     "moderator",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "invalid_role",
			email:    "user@example.com",
			password: "securePassword123",
			role:     "admin",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "empty_email",
			email:    "",
			password: "securePassword123",
			role:     "moderator",
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокаем функцию хеширования пароля для теста "password_hashing_error"

			got, gotErr := domain.NewUser(tt.email, tt.password, tt.role)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if gotErr == nil && tt.want != nil {
				// Если ошибка нет, проверим создание пользователя
				if got.Email != tt.want.Email || got.PasswordHash == "" ||
					got.Role != tt.want.Role {
					t.Errorf("NewUser() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestUser_CheckPasswordHash(t *testing.T) {
	tests := []struct {
		name      string
		cemail    string
		cpassword string
		crole     string
		password  string
		want      bool
	}{
		{
			name:      "valid_password",
			cemail:    "user@example.com",
			cpassword: "securePassword123",
			crole:     "moderator",
			password:  "securePassword123",
			want:      true,
		},
		{
			name:      "invalid_password",
			cemail:    "user@example.com",
			cpassword: "securePassword123",
			crole:     "moderator",
			password:  "wrongPassword",
			want:      false,
		},
		{
			name:      "empty_password",
			cemail:    "user@example.com",
			cpassword: "securePassword123",
			crole:     "moderator",
			password:  "",
			want:      false,
		},
		{
			name:      "password_with_special_characters",
			cemail:    "user@example.com",
			cpassword: "P@ssw0rd123!",
			crole:     "moderator",
			password:  "P@ssw0rd123!",
			want:      true,
		},
		{
			name:      "password_case_sensitive",
			cemail:    "user@example.com",
			cpassword: "SecurePassword",
			crole:     "moderator",
			password:  "securepassword",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := domain.NewUser(tt.cemail, tt.cpassword, tt.crole)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			got := u.CheckPasswordHash(tt.password)
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
