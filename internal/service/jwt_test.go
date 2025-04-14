package service_test

import (
	"testing"
	"time"

	"avito_pvz/internal/service"
)

func TestJWTManager_ValidateToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		secretKey   string
		expiry      time.Duration
		userUUID    string
		role        string
		tokenString string
		want        *service.UserClaims
		wantErr     bool
	}{
		{
			name:      "valid-token",
			secretKey: "test-secret-key",
			expiry:    time.Hour,
			userUUID:  "test-uuid",
			role:      "admin",
			want: &service.UserClaims{
				UUID: "test-uuid",
				Role: "admin",
			},
			wantErr: false,
		},
		{
			name:        "invalid-token",
			secretKey:   "test-secret-key",
			expiry:      time.Hour,
			tokenString: "invalid-token",
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "empty-token",
			secretKey:   "test-secret-key",
			expiry:      time.Hour,
			tokenString: "",
			want:        nil,
			wantErr:     true,
		},
		{
			name:      "expired-token",
			secretKey: "test-secret-key",
			expiry:    -time.Hour,
			userUUID:  "test-uuid",
			role:      "admin",
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := service.NewJWTManager(tt.secretKey, tt.expiry)

			// Генерируем токен только для тестов с валидным и просроченным токеном
			if tt.name == "valid-token" || tt.name == "expired-token" {
				var err error
				tt.tokenString, err = m.GenerateToken(tt.userUUID, tt.role)
				if err != nil {
					t.Fatalf("Не удалось сгенерировать тестовый токен: %v", err)
				}
			}

			gotUUID, gotROLE, err := m.ValidateToken(tt.tokenString)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && gotUUID != "" && gotROLE != "" {
				if gotUUID != tt.want.UUID {
					t.Errorf("ValidateToken() UUID = %v, want %v", gotUUID, tt.want.UUID)
				}

				if gotROLE != tt.want.Role {
					t.Errorf("ValidateToken() Role = %v, want %v", gotROLE, tt.want.Role)
				}
			}
		})
	}
}

func TestJWTManager_GenerateToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		secretKey string
		expiry    time.Duration
		userUUID  string
		role      string
		wantErr   bool
	}{
		{
			name:      "successful-token-generation",
			secretKey: "test-secret-key",
			expiry:    time.Hour,
			userUUID:  "test-uuid",
			role:      "admin",
			wantErr:   false,
		},
		{
			name:      "empty-uuid",
			secretKey: "test-secret-key",
			expiry:    time.Hour,
			userUUID:  "",
			role:      "admin",
			wantErr:   false,
		},
		{
			name:      "empty-role",
			secretKey: "test-secret-key",
			expiry:    time.Hour,
			userUUID:  "test-uuid",
			role:      "",
			wantErr:   false,
		},
		{
			name:      "empty-secret-key",
			secretKey: "",
			expiry:    time.Hour,
			userUUID:  "test-uuid",
			role:      "admin",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := service.NewJWTManager(tt.secretKey, tt.expiry)
			got, err := m.GenerateToken(tt.userUUID, tt.role)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == "" {
					t.Error("GenerateToken() вернул пустой токен")
				}

				uuid, role, err := m.ValidateToken(got)
				if err != nil {
					t.Errorf("Не удалось валидировать сгенерированный токен: %v", err)
					return
				}

				if uuid != tt.userUUID {
					t.Errorf(
						"GenerateToken() UUID в токене = %v, want %v",
						uuid,
						tt.userUUID,
					)
				}

				if role != tt.role {
					t.Errorf("GenerateToken() Role в токене = %v, want %v", role, tt.role)
				}
			}
		})
	}
}
