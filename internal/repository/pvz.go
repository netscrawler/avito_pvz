package repository

import (
	"context"

	"avito_pvz/internal/models/domain"

	"github.com/google/uuid"
)

type PVZRepository interface {
	Create(ctx context.Context, pvz *domain.PVZ) error
	GetAll(ctx context.Context) ([]domain.PVZ, error)
	GetWithParam(ctx context.Context, params domain.Params) ([]domain.PVZAgregate, error)
	Exist(ctx context.Context, pvz uuid.UUID) error
}

type PVZ struct {
	PVZRepository
}

func NewPVZ(p PVZRepository) *PVZ {
	return &PVZ{
		PVZRepository: p,
	}
}
