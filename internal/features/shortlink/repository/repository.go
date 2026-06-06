package shortlink_reposiotry

import core_pgx "api/internal/core/repository/pgx"

type ShortlinkRepository struct {
	pool core_pgx.Pool
}

func NewShortLinkRepository(
	pool core_pgx.Pool,
) *ShortlinkRepository {
	return &ShortlinkRepository{
		pool: pool,
	}
}
