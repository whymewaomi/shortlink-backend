package shortlink_reposiotry

import (
	"context"
)

func (r *ShortlinkRepository) UpdateCountLink(
	ctx context.Context,
	userID int,
	count int,
) error {
	sql := `
	UPDATE shortlink.shortlink
	SET count = count + $1
	WHERE user_id = $2
	`

	if _, err := r.pool.Exec(ctx, sql, count, userID); err != nil {
		return err
	}

	return nil
}