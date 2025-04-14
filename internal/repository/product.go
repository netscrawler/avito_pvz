package repository

import (
	"context"

	"avito_pvz/internal/models/domain"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetLast(ctx context.Context, receptionID uuid.UUID) (*domain.Product, error)
	Delete(ctx context.Context, product *domain.Product) error
}

type Product struct {
	ProductRepository
}

func NewProduct(p ProductRepository) *Product {
	return &Product{
		ProductRepository: p,
	}
}
