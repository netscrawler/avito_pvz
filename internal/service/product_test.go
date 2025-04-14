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
)

func TestProduct_Create(t *testing.T) {
	tests := []struct {
		name        string
		product     domain.ProductToAdd
		setupMocks  func(*service.MockProductProvider, *service.MockReceptionGetter, *service.MockPVZChecker)
		expected    *domain.Product
		expectedErr error
	}{
		{
			name: "successful product creation",
			product: domain.ProductToAdd{
				UUID: domain.PVZID(uuid.Max),
				Type: domain.ProductType(domain.ProductTypeElectronics),
			},
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				reception := &domain.Reception{
					ID:        uuid.Max,
					PvzID:     pvzID,
					Status:    domain.ReceptionStatusInProgress,
					CreatedAt: time.Now(),
				}

				mc.On("Exist", mock.Anything, pvzID).Return(nil)
				mr.On("GetLast", mock.Anything, pvzID).Return(reception, nil)
				mp.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expected: &domain.Product{
				Type: domain.ProductType(domain.ProductTypeElectronics),
			},
			expectedErr: nil,
		},
		{
			name: "invalid product type",
			product: domain.ProductToAdd{
				UUID: domain.PVZID(uuid.Max),
				Type: "invalid_type",
			},
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
			},
			expected:    nil,
			expectedErr: models.ErrInvalidProductType,
		},
		{
			name: "PVZ not found",
			product: domain.ProductToAdd{
				UUID: domain.PVZID(uuid.Max),
				Type: domain.ProductType(domain.ProductTypeElectronics),
			},
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				mc.On("Exist", mock.Anything, pvzID).Return(domain.ErrPVZNotExist)
			},
			expected:    nil,
			expectedErr: models.ErrPVZNotFound,
		},
		{
			name: "reception not found",
			product: domain.ProductToAdd{
				UUID: domain.PVZID(uuid.Max),
				Type: domain.ProductType(domain.ProductTypeElectronics),
			},
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				mc.On("Exist", mock.Anything, pvzID).Return(nil)
				mr.On("GetLast", mock.Anything, pvzID).Return(nil, domain.ErrNotFound)
			},
			expected:    nil,
			expectedErr: models.ErrReceptionDontExist,
		},
		{
			name: "reception already closed",
			product: domain.ProductToAdd{
				UUID: domain.PVZID(uuid.Max),
				Type: domain.ProductType(domain.ProductTypeElectronics),
			},
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				reception := &domain.Reception{
					ID:        uuid.Max,
					PvzID:     pvzID,
					Status:    domain.ReceptionStatusClosed,
					CreatedAt: time.Now(),
				}

				mc.On("Exist", mock.Anything, pvzID).Return(nil)
				mr.On("GetLast", mock.Anything, pvzID).Return(reception, nil)
			},
			expected:    nil,
			expectedErr: models.ErrReceptionAlreadyClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := service.NewMockProductProvider(t)
			mockReception := service.NewMockReceptionGetter(t)
			mockPVZ := service.NewMockPVZChecker(t)

			tt.setupMocks(mockProduct, mockReception, mockPVZ)

			// Create service
			service := service.NewProduct(mockProduct, mockReception, mockPVZ)

			// Call method
			result, err := service.Create(context.Background(), tt.product)

			// Assert results
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Type, result.Type)
			}

			// Verify service
			mockProduct.AssertExpectations(t)
			mockReception.AssertExpectations(t)
			mockPVZ.AssertExpectations(t)
		})
	}
}

func TestProduct_DeleteLast(t *testing.T) {
	tests := []struct {
		name        string
		pvzID       domain.PVZID
		setupMocks  func(*service.MockProductProvider, *service.MockReceptionGetter, *service.MockPVZChecker)
		expectedErr error
	}{
		{
			name:  "successful product deletion",
			pvzID: domain.PVZID(uuid.Max),
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				reception := &domain.Reception{
					ID:        uuid.Max,
					PvzID:     pvzID,
					Status:    domain.ReceptionStatusInProgress,
					CreatedAt: time.Now(),
				}
				product := &domain.Product{
					ID:          uuid.Max,
					ReceptionID: reception.ID,
					Type:        domain.ProductTypeElectronics,
					CreatedAt:   time.Now(),
				}

				mc.On("Exist", mock.Anything, pvzID).Return(nil)
				mr.On("GetLast", mock.Anything, pvzID).Return(reception, nil)
				mp.On("GetLast", mock.Anything, reception.ID).Return(product, nil)
				mp.On("Delete", mock.Anything, product).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "PVZ not found",
			pvzID: domain.PVZID(uuid.Max),
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				mc.On("Exist", mock.Anything, pvzID).Return(domain.ErrPVZNotExist)
			},
			expectedErr: models.ErrPVZNotFound,
		},
		{
			name:  "reception not found",
			pvzID: domain.PVZID(uuid.Max),
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				mc.On("Exist", mock.Anything, pvzID).Return(nil)
				mr.On("GetLast", mock.Anything, pvzID).Return(nil, domain.ErrNotFound)
			},
			expectedErr: models.ErrReceptionDontExist,
		},
		{
			name:  "reception already closed",
			pvzID: domain.PVZID(uuid.Max),
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				reception := &domain.Reception{
					ID:        uuid.New(),
					PvzID:     pvzID,
					Status:    domain.ReceptionStatusClosed,
					CreatedAt: time.Now(),
				}

				mc.On("Exist", mock.Anything, pvzID).Return(nil)
				mr.On("GetLast", mock.Anything, pvzID).Return(reception, nil)
			},
			expectedErr: models.ErrReceptionAlreadyClosed,
		},
		{
			name:  "product not found",
			pvzID: domain.PVZID(uuid.Max),
			setupMocks: func(mp *service.MockProductProvider, mr *service.MockReceptionGetter, mc *service.MockPVZChecker) {
				pvzID := uuid.Max
				reception := &domain.Reception{
					ID:        uuid.Max,
					PvzID:     pvzID,
					Status:    domain.ReceptionStatusInProgress,
					CreatedAt: time.Now(),
				}

				mc.On("Exist", mock.Anything, pvzID).Return(nil)
				mr.On("GetLast", mock.Anything, pvzID).Return(reception, nil)
				mp.On("GetLast", mock.Anything, reception.ID).Return(nil, domain.ErrNotFound)
			},
			expectedErr: models.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service
			mockProduct := service.NewMockProductProvider(t)
			mockReception := service.NewMockReceptionGetter(t)
			mockPVZ := service.NewMockPVZChecker(t)

			// Setup service
			tt.setupMocks(mockProduct, mockReception, mockPVZ)

			// Create service
			service := service.NewProduct(mockProduct, mockReception, mockPVZ)

			// Call method
			err := service.DeleteLast(context.Background(), tt.pvzID)

			// Assert results
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Verify service
			mockProduct.AssertExpectations(t)
			mockReception.AssertExpectations(t)
			mockPVZ.AssertExpectations(t)
		})
	}
}
