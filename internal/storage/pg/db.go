package postgres

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	log     *slog.Logger
	DB      *pgxpool.Pool
	Builder *squirrel.StatementBuilderType
}

func MustSetup(ctx context.Context, dsn string, log *slog.Logger) *Storage {
	//nolint: varnamelen
	const op = "storage.postgres.Setup"

	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	pgConn, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		panic(err)
	}

	err = pgConn.Ping(ctx)
	if err != nil {
		panic(err)
	}

	log.Info(op + "Successfyly connect to database")

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Storage{
		log:     log,
		DB:      pgConn,
		Builder: &builder,
	}
}

func (s *Storage) Stop() {
	const op = "storage.pg.Stop"

	s.DB.Close()

	s.log.Info(op + "Connection to database closed")
}
