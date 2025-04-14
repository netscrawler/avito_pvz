package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"avito_pvz/internal/models/domain"
	postgres "avito_pvz/internal/storage/pg"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type pgReception struct {
	storage *postgres.Storage
}

func NewPgReception(db *postgres.Storage) *pgReception {
	return &pgReception{
		storage: db,
	}
}

func (p *pgReception) Close(ctx context.Context, reception domain.Reception) error {
	query, args, err := p.storage.Builder.
		Update("receptions").
		Set("status", reception.Status).
		Where(squirrel.Eq{"id": reception.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	ct, err := p.storage.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (p *pgReception) GetLast(ctx context.Context, pvz uuid.UUID) (*domain.Reception, error) {
	query, args, err := p.storage.Builder.
		Select("id", "pvz_id", "status", "created_at").
		From("receptions").
		Where(squirrel.Eq{"pvz_id": pvz}).
		OrderBy("created_at DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	row := p.storage.DB.QueryRow(ctx, query, args...)

	var reception domain.Reception
	if err := row.Scan(&reception.ID, &reception.PvzID, &reception.Status, &reception.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	return &reception, nil
}

func (p *pgReception) Create(ctx context.Context, reception domain.Reception) error {
	query, args, err := p.storage.Builder.
		Insert("receptions").
		Columns("pvz_id", "status").
		Values(reception.PvzID, reception.Status).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	err = p.storage.DB.QueryRow(ctx, query, args...).Scan(&reception.ID, &reception.CreatedAt)
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	return nil
}
