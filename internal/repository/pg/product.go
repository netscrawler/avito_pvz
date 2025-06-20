package pgrepo

import (
	"avito_pvz/internal/models/domain"
	"context"
	"errors"
	"fmt"

	postgres "avito_pvz/internal/storage/pg"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type pgProduct struct {
	db *postgres.Storage
}

func NewPgProduct(db *postgres.Storage) *pgProduct {
	return &pgProduct{
		db: db,
	}
}

func (p *pgProduct) Create(ctx context.Context, product *domain.Product) error {
	query, args, err := squirrel.
		Insert("products").
		Columns("reception_id", "product_type").
		Values(product.ReceptionID, product.Type).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	row := p.db.DB.QueryRow(ctx, query, args...)
	if err := row.Scan(&product.ID, &product.CreatedAt); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return nil
}

func (p *pgProduct) GetLast(ctx context.Context, receptionID uuid.UUID) (*domain.Product, error) {
	query, args, err := squirrel.
		Select("id", "reception_id", "product_type", "created_at").
		From("products").
		Where(squirrel.Eq{"reception_id": receptionID}).
		OrderBy("created_at DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	row := p.db.DB.QueryRow(ctx, query, args...)

	var product domain.Product
	if err := row.Scan(&product.ID, &product.ReceptionID, &product.Type, &product.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return &product, nil
}

func (p *pgProduct) Delete(ctx context.Context, product *domain.Product) error {
	query, args, err := squirrel.
		Delete("products").
		Where(squirrel.Eq{"id": product.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return nil
}
