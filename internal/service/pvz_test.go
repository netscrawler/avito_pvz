package service_test

import (
	"context"
	"testing"
	"time"

	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"
	"avito_pvz/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPVZ_List(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		setupMocks func(*service.MockPVZProvider)
		want       *domain.PVZList
		wantErr    error
	}{
		{
			name: "no pvz found",
			setupMocks: func(mp *service.MockPVZProvider) {
				mp.On("GetWithParam", mock.Anything, mock.Anything).
					Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: models.ErrPVZNotFound,
		},
		{
			name: "repo returns internal error",
			setupMocks: func(mp *service.MockPVZProvider) {
				mp.On("GetWithParam", mock.Anything, mock.Anything).
					Return(nil, domain.ErrInternal)
			},
			want:    nil,
			wantErr: models.ErrInternal,
		},
		{
			name: "successful list",
			setupMocks: func(mp *service.MockPVZProvider) {
				mockDomainPVZs := []domain.PVZ{
					{
						ID:               (*domain.PVZID)(&uuid.Max),
						City:             "Москва",
						RegistrationDate: time.Time{},
					},
				}
				mp.On("GetWithParam", mock.Anything, mock.Anything).
					Return(mockDomainPVZs, nil)
			},
			want: &domain.PVZList{
				{
					ID:               (*domain.PVZID)(&uuid.Max),
					City:             "Москва",
					RegistrationDate: time.Time{},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZ := service.NewMockPVZProvider(t)
			tt.setupMocks(mockPVZ)
			service := service.NewPVZServce(mockPVZ)

			params := domain.Params{}
			got, err := service.List(context.Background(), params)

			// Assert results
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want, got)
			}

			mockPVZ.AssertExpectations(t)
		})
	}
}

func TestPVZ_GetAllPVZ(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		setupMocks func(*service.MockPVZProvider)
		want       domain.PVZList
		wantErr    error
	}{
		{
			name: "successful get pvz list",
			setupMocks: func(mp *service.MockPVZProvider) {
				pvzs := []domain.PVZ{{
					ID:               (*domain.PVZID)(&uuid.Max),
					City:             "Москва",
					RegistrationDate: time.Time{},
				}}
				mp.On("GetAll", mock.Anything).
					Return(pvzs, nil)
			},
			want: domain.PVZList{{
				ID:               (*domain.PVZID)(&uuid.Max),
				RegistrationDate: time.Time{},
				City:             "Москва",
			}},
			wantErr: nil,
		},
		{
			name: "empty pvz list",
			setupMocks: func(mp *service.MockPVZProvider) {
				mp.On("GetAll", mock.Anything).
					Return([]domain.PVZ{}, nil)
			},
			want:    domain.PVZList{},
			wantErr: nil,
		},
		{
			name: "pvz not found error from repo",
			setupMocks: func(mp *service.MockPVZProvider) {
				mp.On("GetAll", mock.Anything).
					Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: models.ErrPVZNotFound,
		},
		{
			name: "unexpected internal error from repo",
			setupMocks: func(mp *service.MockPVZProvider) {
				mp.On("GetAll", mock.Anything).
					Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: models.ErrInternal,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZ := service.NewMockPVZProvider(t)
			tt.setupMocks(mockPVZ)
			service := service.NewPVZServce(mockPVZ)

			got, err := service.GetAllPVZ(context.Background())

			// Assert results
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want, got)
			}
			mockPVZ.AssertExpectations(t)
		})
	}
}

func TestPVZ_Create(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		pvz        domain.PvzCity
		setupMocks func(*service.MockPVZProvider)
		want       *domain.PVZ
		wantErr    error
	}{
		{
			name: "successful create pvz",
			pvz:  "Москва",
			setupMocks: func(mp *service.MockPVZProvider) {
				mp.On("Create", mock.Anything, mock.Anything).
					Return(nil)
			},
			want: &domain.PVZ{
				ID:               (*domain.PVZID)(&uuid.Max),
				City:             "Москва",
				RegistrationDate: time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "invalid city",
			pvz:  "!!!",
			setupMocks: func(mp *service.MockPVZProvider) {
			},
			want:    nil,
			wantErr: models.ErrInvalidCity,
		},
		{
			name: "repo returns internal error",
			pvz:  "Казань",
			setupMocks: func(mp *service.MockPVZProvider) {
				mp.On("Create", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
			want:    nil,
			wantErr: models.ErrInternal,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZ := service.NewMockPVZProvider(t)
			tt.setupMocks(mockPVZ)
			service := service.NewPVZServce(mockPVZ)

			got, err := service.Create(context.Background(), tt.pvz)

			// Assert results
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)

				require.Equal(t, tt.pvz, got.City)

				require.NotNil(t, got.ID)
				require.NotEqual(t, uuid.Nil, *got.ID) // Убедись, что ID валидный

				require.NotNil(t, got.RegistrationDate)
				require.WithinDuration(t, time.Now(), got.RegistrationDate, time.Second)

			}
			mockPVZ.AssertExpectations(t)
		})
	}
}
