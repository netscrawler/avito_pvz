package service_test

import (
	"context"
	"errors"
	"reflect"
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

func TestReception_CloseLastReception(t *testing.T) {
	id := domain.PVZID(uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	activeReception := &domain.Reception{
		ID:        uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		PvzID:     uuid.UUID(id),
		CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Status:    "in_progress",
	}
	inactiveReception := &domain.Reception{
		ID:        uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		PvzID:     uuid.UUID(id),
		CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Status:    "close",
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		setupMocks func(*service.MockPVZChecker, *service.MockReceptionProvider)
		wantErr    error
	}{
		{
			name: "pvz not found",
			setupMocks: func(pvz *service.MockPVZChecker, rp *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.UUID(id)).
					Return(domain.ErrNotFound)
			},
			wantErr: models.ErrPVZNotFound,
		},
		{
			name: "pvz exist returns internal error",
			setupMocks: func(pvz *service.MockPVZChecker, rp *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.UUID(id)).
					Return(errors.New("db fail"))
			},
			wantErr: models.ErrInternal,
		},
		{
			name: "reception not found",
			setupMocks: func(pvz *service.MockPVZChecker, rp *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.UUID(id)).
					Return(nil)
				rp.On("GetLast", mock.Anything, uuid.UUID(id)).
					Return(nil, domain.ErrNotFound)
			},
			wantErr: models.ErrReceptionDontExist,
		},
		{
			name: "reception already closed",
			setupMocks: func(pvz *service.MockPVZChecker, rp *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.UUID(id)).
					Return(nil)
				rp.On("GetLast", mock.Anything, uuid.UUID(id)).
					Return(inactiveReception, nil)
			},
			wantErr: models.ErrReceptionAlreadyClosed,
		},
		{
			name: "reception get returns internal error",
			setupMocks: func(pvz *service.MockPVZChecker, rp *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.UUID(id)).
					Return(nil)
				rp.On("GetLast", mock.Anything, uuid.UUID(id)).
					Return(nil, errors.New("db error"))
			},
			wantErr: models.ErrInternal,
		},
		{
			name: "close fails",
			setupMocks: func(pvz *service.MockPVZChecker, rp *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.UUID(id)).
					Return(nil)
				rp.On("GetLast", mock.Anything, uuid.UUID(id)).
					Return(activeReception, nil)
				rp.On("Close", mock.Anything, *activeReception).
					Return(errors.New("db error"))
			},
			wantErr: models.ErrInternal,
		},
		{
			name: "successful close",
			setupMocks: func(pvz *service.MockPVZChecker, rp *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.UUID(id)).
					Return(nil)
				rp.On("GetLast", mock.Anything, uuid.UUID(id)).
					Return(activeReception, nil)
				rp.On("Close", mock.Anything, *activeReception).
					Return(nil)
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZ := service.NewMockPVZChecker(t)
			mockReception := service.NewMockReceptionProvider(t)
			tt.setupMocks(mockPVZ, mockReception)

			svc := service.NewReceptionService(mockReception, mockPVZ)

			_, err := svc.CloseLastReception(context.Background(), id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			mockPVZ.AssertExpectations(t)
			mockReception.AssertExpectations(t)
		})
	}
}

func TestReception_Create(t *testing.T) {
	tests := []struct {
		name       string
		pvzID      domain.PVZID
		want       *domain.Reception
		wantErr    error
		setupMocks func(pvz *service.MockPVZChecker, reception *service.MockReceptionProvider)
	}{
		{
			name:  "reception found successfully",
			pvzID: domain.PVZID(uuid.Max),
			want: &domain.Reception{
				ID:        uuid.Max,
				CreatedAt: time.Time{},
				PvzID:     uuid.Max,
				Status:    "in_progress",
			},
			wantErr: nil,
			setupMocks: func(pvz *service.MockPVZChecker, reception *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.Max).Return(nil)
				reception.On("GetLast", mock.Anything, uuid.Max).Return(
					&domain.Reception{
						ID:        uuid.Max,
						PvzID:     uuid.Max,
						Status:    "close",
						CreatedAt: time.Time{},
					}, nil,
				)
				reception.On("Create", mock.Anything, mock.MatchedBy(func(r domain.Reception) bool {
					return r.PvzID == uuid.Max &&
						r.Status == domain.ReceptionStatusInProgress
				})).Return(nil)
			},
		},
		{
			name:    "reception already exist",
			pvzID:   domain.PVZID(uuid.Max),
			want:    nil,
			wantErr: models.ErrReceptionAlreadyExist,
			setupMocks: func(pvz *service.MockPVZChecker, reception *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.Max).Return(nil)
				reception.On("GetLast", mock.Anything, uuid.Max).Return(
					&domain.Reception{
						ID:        uuid.Max,
						PvzID:     uuid.Max,
						Status:    "in_progress",
						CreatedAt: time.Time{},
					}, nil,
				)
			},
		},
		{
			name:    "pvz not found",
			pvzID:   domain.PVZID(uuid.Max),
			want:    nil,
			wantErr: models.ErrPVZNotFound,
			setupMocks: func(pvz *service.MockPVZChecker, reception *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.Max).Return(domain.ErrNotFound)
			},
		},
		{
			name:    "Exist fails internally",
			pvzID:   domain.PVZID(uuid.Max),
			want:    nil,
			wantErr: models.ErrInternal,
			setupMocks: func(pvz *service.MockPVZChecker, reception *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.Max).Return(assert.AnError)
			},
		},
		{
			name:    "GetLast fails internally",
			pvzID:   domain.PVZID(uuid.Max),
			want:    nil,
			wantErr: models.ErrInternal,
			setupMocks: func(pvz *service.MockPVZChecker, reception *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.Max).Return(nil)
				reception.On("GetLast", mock.Anything, uuid.Max).Return(nil, assert.AnError)
			},
		},
		{
			name:    "Create fails",
			pvzID:   domain.PVZID(uuid.Max),
			want:    nil,
			wantErr: models.ErrInternal,
			setupMocks: func(pvz *service.MockPVZChecker, reception *service.MockReceptionProvider) {
				pvz.On("Exist", mock.Anything, uuid.Max).Return(nil)
				reception.On("GetLast", mock.Anything, uuid.Max).Return(
					&domain.Reception{
						ID:        uuid.Max,
						PvzID:     uuid.Max,
						Status:    "close",
						CreatedAt: time.Time{},
					}, nil,
				)
				reception.On("Create", mock.Anything, mock.MatchedBy(func(r domain.Reception) bool {
					return r.PvzID == uuid.Max &&
						r.Status == domain.ReceptionStatusInProgress
				})).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZ := service.NewMockPVZChecker(t)
			mockReception := service.NewMockReceptionProvider(t)
			tt.setupMocks(mockPVZ, mockReception)

			svc := service.NewReceptionService(mockReception, mockPVZ)

			got, err := svc.Create(context.Background(), tt.pvzID)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.want.PvzID, got.PvzID)
				require.WithinDuration(t, time.Now(), got.CreatedAt, time.Second)
				require.Equal(t, reflect.TypeOf(tt.want.ID), reflect.TypeOf(got.ID))
			}

			mockPVZ.AssertExpectations(t)
			mockReception.AssertExpectations(t)
		})
	}
}
