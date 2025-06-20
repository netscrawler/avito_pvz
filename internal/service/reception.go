package service

import (
	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

type ReceptionProvider interface {
	Close(ctx context.Context, reception domain.Reception) error
	GetLast(ctx context.Context, pvz uuid.UUID) (*domain.Reception, error)
	Create(ctx context.Context, reception domain.Reception) error
}
type Reception struct {
	reception ReceptionProvider
	pvz       PVZChecker
}

func (r *Reception) CloseLastReception(
	ctx context.Context,
	pvzID domain.PVZID,
) (*domain.Reception, error) {
	err := r.pvz.Exist(ctx, uuid.UUID(pvzID))
	if errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrPVZNotFound
	}

	if err != nil {
		return nil, models.ErrInternal
	}

	reception, err := r.reception.GetLast(ctx, uuid.UUID(pvzID))
	if reception != nil && !reception.IsActive() {
		return nil, models.ErrReceptionAlreadyClosed
	}

	if errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrReceptionDontExist
	}

	if err != nil {
		return nil, models.ErrInternal
	}

	err = r.reception.Close(ctx, *reception)
	if err != nil {
		return nil, models.ErrInternal
	}

	reception.Close()

	return reception, nil
}

func (r *Reception) Create(ctx context.Context, pvzID domain.PVZID) (*domain.Reception, error) {
	err := r.pvz.Exist(ctx, uuid.UUID(pvzID))
	if errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrPVZNotFound
	}

	if err != nil {
		return nil, models.ErrInternal
	}

	oldReception, err := r.reception.GetLast(ctx, uuid.UUID(pvzID))
	if oldReception != nil && oldReception.IsActive() {
		return nil, models.ErrReceptionAlreadyExist
	}

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrInternal
	}

	reception := domain.NewReception(uuid.UUID(pvzID))

	err = r.reception.Create(ctx, *reception)
	if err != nil {
		return nil, models.ErrInternal
	}

	return reception, nil
}

func NewReceptionService(reception ReceptionProvider, pvz PVZChecker) *Reception {
	return &Reception{
		reception: reception,
		pvz:       pvz,
	}
}
