package service

import (
	"context"
	"errors"

	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"
)

type PVZProvider interface {
	Create(ctx context.Context, pvz *domain.PVZ) error
	GetAll(ctx context.Context) ([]domain.PVZ, error)
	GetWithParam(ctx context.Context, params domain.Params) ([]domain.PVZAgregate, error)
}

type PVZ struct {
	repo PVZProvider
}

func (p *PVZ) GetAllPVZ(ctx context.Context) (domain.PVZList, error) {
	pvzs, err := p.repo.GetAll(ctx)

	if errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrPVZNotFound
	}

	if err != nil {
		return nil, models.ErrInternal
	}

	out := make(domain.PVZList, 0, len(pvzs))
	for _, v := range pvzs {
		out = append(out, &v)
	}
	return out, nil
}

func (p *PVZ) List(ctx context.Context, params domain.Params) (*[]domain.PVZAgregate, error) {
	pvzs, err := p.repo.GetWithParam(ctx, params)

	if errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrPVZNotFound
	}

	if err != nil {
		return nil, models.ErrInternal
	}

	return &pvzs, nil
}

func (p *PVZ) Create(ctx context.Context, city domain.PvzCity) (*domain.PVZ, error) {
	valCity := domain.PvzCity(city)
	if !valCity.IsValid() {
		return nil, models.ErrInvalidCity
	}

	pvz := domain.NewPVZ(valCity)

	err := p.repo.Create(ctx, pvz)
	if err != nil {
		return nil, models.ErrInternal
	}

	return pvz, nil
}

func NewPVZServce(repo PVZProvider) *PVZ {
	return &PVZ{
		repo: repo,
	}
}
