package core_pgx

import (
	core_config "api/internal/core/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

func NewConnect(
	ctx context.Context,
	cfg *core_config.Config,
) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, fmt.Sprintf(
		"postgres://%s:%s@postgres:5432/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
	))
}
