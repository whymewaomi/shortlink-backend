package shortlink_reposiotry

import (
	"api/internal/core/domain"
	"context"
	"fmt"
)

func (r *ShortlinkRepository) CreateUser(
	ctx context.Context,
	user *domain.User,
) (int, error) {
	sql := `
    INSERT INTO shortlink.user (username, email, password_hash)
    VALUES ($1, $2, $3)
    RETURNING id
   `

	var userID int
	if err := r.pool.QueryRow(
		ctx,
		sql,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&userID); err != nil {
		return 0, fmt.Errorf("exec error: %w", err)
	}

	return userID, nil
}

func (r *ShortlinkRepository) CreateShortlink(
	ctx context.Context,
	sl *domain.Shortlink,
) error {
	sql := `
    INSERT INTO shortlink.shortlink (user_id, shortlink, original_link)
    VALUES ($1, $2, $3)
   `

	if _, err := r.pool.Exec(
		ctx,
		sql,
		sl.UserID,
		sl.Shortlink,
		sl.OriginalUrl,
	); err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}

func (r *ShortlinkRepository) CreateShortLinkActivate(
	ctx context.Context,
	sli *domain.ShortlinkInfo,
) error {
	sql := `
    INSERT INTO shortlink.shortlink_activate (ip_addr, user_agent, shortlink_id)
    VALUES ($1, $2, $3)
   `

	if _, err := r.pool.Exec(
		ctx,
		sql,
		sli.IPAddr,
		sli.UserAgent,
		sli.ShortlinkID,
	); err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}
