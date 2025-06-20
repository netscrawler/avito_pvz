package repository

import (
	"avito_pvz/internal/models/domain"
	"context"

	"github.com/google/uuid"
)

type ReceptionRepository interface {
	Close(ctx context.Context, reception domain.Reception) error
	GetLast(ctx context.Context, pvz uuid.UUID) (*domain.Reception, error)
	Create(ctx context.Context, reception domain.Reception) error
}

type Reception struct {
	ReceptionRepository
}

func NewReception(r ReceptionRepository) *Reception {
	return &Reception{
		ReceptionRepository: r,
	}
}
