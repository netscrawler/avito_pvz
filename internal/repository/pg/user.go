package pgrepo

import (
	"avito_pvz/internal/models/domain"
	"context"
	"errors"
	"fmt"

	postgres "avito_pvz/internal/storage/pg"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgUser struct {
	storage *postgres.Storage
}

func NewPgUser(db *postgres.Storage) *pgUser {
	return &pgUser{
		storage: db,
	}
}

func (p *pgUser) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query, args, err := p.storage.Builder.
		Select("id", "email", "password_hash", "role", "created_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	row := p.storage.DB.QueryRow(ctx, query, args...)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	return &user, nil
}

func (p *pgUser) Create(ctx context.Context, user *domain.User) error {
	query, args, err := p.storage.Builder.
		Insert("users").
		Columns("email", "password_hash", "role").
		Values(user.Email, user.PasswordHash, user.Role).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	err = p.storage.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	return nil
}
