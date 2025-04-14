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

type pgPvz struct {
	storage *postgres.Storage
}

func NewPgPvz(db *postgres.Storage) *pgPvz {
	return &pgPvz{
		storage: db,
	}
}

func (p *pgPvz) Create(ctx context.Context, pvz *domain.PVZ) error {
	query, args, err := p.storage.Builder.
		Insert("pvzs").
		Columns("id", "city", "created_at").
		Values(pvz.ID, pvz.City, pvz.RegistrationDate).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	_, err = p.storage.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	return nil
}

func (p *pgPvz) GetAll(ctx context.Context) ([]domain.PVZ, error) {
	query, args, err := p.storage.Builder.
		Select("id", "city", "created_at").
		From("pvzs").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	rows, err := p.storage.DB.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	list := make([]domain.PVZ, 0, 1)

	for rows.Next() {
		var pvz domain.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate); err != nil {
			// TODO: add log
			continue
		}
		list = append(list, pvz)
	}

	return list, nil
}

// func (p *pgPvz) GetWithParam(ctx context.Context, params domain.Params) ([]domain.PVZ, error) {
// 	qb := p.storage.Builder.
// 		Select("id", "city", "created_at").
// 		From("pvzs")
//
// 	if params.StartDate != nil {
// 		qb = qb.Where(squirrel.GtOrEq{"created_at": *params.StartDate})
// 	}
// 	if params.EndDate != nil {
// 		qb = qb.Where(squirrel.LtOrEq{"created_at": *params.EndDate})
// 	}
//
// 	if params.Limit != nil {
// 		qb = qb.Limit(uint64(*params.Limit))
// 	}
// 	if params.Page != nil && params.Limit != nil {
// 		offset := (*params.Page - 1) * (*params.Limit)
// 		qb = qb.Offset(uint64(offset))
// 	}
//
// 	query, args, err := qb.ToSql()
// 	if err != nil {
// 		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
// 	}
//
// 	rows, err := p.storage.DB.Query(ctx, query, args...)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, domain.ErrNotFound
// 		}
// 		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
// 	}
// 	defer rows.Close()
//
// 	list := make([]domain.PVZ, 0, 1)
// 	for rows.Next() {
// 		var pvz domain.PVZ
// 		if err := rows.Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate); err != nil {
// 			// TODO: log or wrap error
// 			continue
// 		}
// 		list = append(list, pvz)
// 	}
//
// 	return list, nil
// }

func (p *pgPvz) Exist(ctx context.Context, pvz uuid.UUID) error {
	query, args, err := p.storage.Builder.Select("id", "city", "created_at").
		From("pvzs").
		Where(squirrel.Eq{"id": pvz}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	row := p.storage.DB.QueryRow(ctx, query, args...)
	var pvztst domain.PVZ
	err = row.Scan(&pvztst.ID, &pvztst.City, &pvztst.RegistrationDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	return nil
}

func (p *pgPvz) GetWithParam(
	ctx context.Context,
	params domain.Params,
) ([]domain.PVZAgregate, error) {
	// Строим запрос для получения данных о ПВЗ
	qb := p.storage.Builder.
		Select("pvzs.id", "pvzs.city", "pvzs.created_at").
		From("pvzs")

	if params.StartDate != nil {
		qb = qb.Where(squirrel.GtOrEq{"pvzs.created_at": *params.StartDate})
	}
	if params.EndDate != nil {
		qb = qb.Where(squirrel.LtOrEq{"pvzs.created_at": *params.EndDate})
	}

	if params.Limit != nil {
		qb = qb.Limit(uint64(*params.Limit))
	}
	if params.Page != nil && params.Limit != nil {
		offset := (*params.Page - 1) * (*params.Limit)
		qb = qb.Offset(uint64(offset))
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	// Выполняем запрос для ПВЗ
	rows, err := p.storage.DB.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}
	defer rows.Close()

	var pvzs []domain.PVZ
	for rows.Next() {
		var pvz domain.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate); err != nil {
			return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
		}
		pvzs = append(pvzs, pvz)
	}

	// Теперь получаем данные о приемках и продуктах для каждого ПВЗ
	var result []domain.PVZAgregate
	for _, pvz := range pvzs {
		// Получаем приемки для каждого ПВЗ
		receptions, err := p.getReceptionsByPVZID(ctx, pvz.ID.String())
		if err != nil {
			return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
		}

		// Для каждой приемки получаем продукты
		var receptionDetails []struct {
			Products  *[]domain.Product
			Reception *domain.Reception
		}
		for _, reception := range receptions {
			products, err := p.getProductsByReceptionID(ctx, reception.ID.String())
			if err != nil {
				return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
			}

			receptionDetails = append(receptionDetails, struct {
				Products  *[]domain.Product
				Reception *domain.Reception
			}{
				Products:  &products,
				Reception: &reception,
			})
		}

		result = append(result, domain.PVZAgregate{
			Pvz:        &pvz,
			Receptions: &receptionDetails,
		})
	}

	return result, nil
}

// getReceptionsByPVZID выполняет запрос для получения приемок по идентификатору ПВЗ с использованием squirrel
func (p *pgPvz) getReceptionsByPVZID(
	ctx context.Context,
	pvzID string,
) ([]domain.Reception, error) {
	qb := p.storage.Builder.
		Select("id", "pvz_id", "status", "created_at").
		From("recepcions").
		Where(squirrel.Eq{"pvz_id": pvzID})

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	rows, err := p.storage.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}
	defer rows.Close()

	var receptions []domain.Reception
	for rows.Next() {
		var reception domain.Reception
		if err := rows.Scan(&reception.ID, &reception.PvzID, &reception.Status, &reception.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
		}
		receptions = append(receptions, reception)
	}

	return receptions, nil
}

// getProductsByReceptionID выполняет запрос для получения продуктов по идентификатору приемки с использованием squirrel
func (p *pgPvz) getProductsByReceptionID(
	ctx context.Context,
	receptionID string,
) ([]domain.Product, error) {
	qb := p.storage.Builder.
		Select("id", "reception_id", "product_type", "created_at").
		From("products").
		Where(squirrel.Eq{"reception_id": receptionID})

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	rows, err := p.storage.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.ReceptionID, &product.Type, &product.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
		}
		products = append(products, product)
	}

	return products, nil
}
