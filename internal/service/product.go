package service

import (
	"context"
	"errors"

	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"

	"github.com/google/uuid"
)

type ProductProvider interface {
	Create(ctx context.Context, product *domain.Product) error
	GetLast(ctx context.Context, receptionID uuid.UUID) (*domain.Product, error)
	Delete(ctx context.Context, product *domain.Product) error
}

type ReceptionGetter interface {
	GetLast(ctx context.Context, pvz uuid.UUID) (*domain.Reception, error)
}

type PVZChecker interface {
	Exist(ctx context.Context, pvz uuid.UUID) error
}

type Product struct {
	product   ProductProvider
	reception ReceptionGetter
	pvz       PVZChecker
}

func (p *Product) Create(
	ctx context.Context,
	product domain.ProductToAdd,
) (*domain.Product, error) {
	pType := domain.ProductType(product.Type)
	if !pType.IsValid() {
		return nil, models.ErrInvalidProductType
	}

	reception, err := p.getActiveReceprion(ctx, uuid.UUID(product.UUID))
	if err != nil {
		return nil, err
	}

	prod := domain.NewProduct(reception.ID, pType)

	err = p.product.Create(ctx, prod)
	if err != nil {
		return nil, models.ErrInternal
	}

	return prod, nil
}

func (p *Product) getActiveReceprion(
	ctx context.Context,
	pvzID uuid.UUID,
) (*domain.Reception, error) {
	err := p.pvz.Exist(ctx, pvzID)
	if errors.Is(err, domain.ErrPVZNotExist) {
		return nil, models.ErrPVZNotFound
	}
	if err != nil {
		return nil, models.ErrInternal
	}

	reception, err := p.reception.GetLast(ctx, pvzID)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, models.ErrReceptionDontExist
	}

	if err != nil {
		return nil, models.ErrInternal
	}

	if !reception.IsActive() {
		return nil, models.ErrReceptionAlreadyClosed
	}

	return reception, nil
}

func (p *Product) DeleteLast(ctx context.Context, pvzID domain.PVZID) error {
	reception, err := p.getActiveReceprion(ctx, uuid.UUID(pvzID))
	if err != nil {
		return err
	}

	product, err := p.product.GetLast(ctx, reception.ID)
	if errors.Is(err, domain.ErrNotFound) {
		return models.ErrProductNotFound
	}

	err = p.product.Delete(ctx, product)
	if err != nil {
		return models.ErrInternal
	}

	return nil
}

func NewProduct(product ProductProvider, reception ReceptionGetter, pvz PVZChecker) *Product {
	return &Product{
		product:   product,
		reception: reception,
		pvz:       pvz,
	}
}
