package shortlink_reposiotry

import (
	"api/internal/core/domain"
	"context"
	"fmt"
)

func (r *ShortlinkRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*domain.User, error) {
	sql := `
   SELECT id, username, email, password_hash
   FROM shortlink.user
   WHERE username = $1
   `

	var u domain.User
	if err := r.pool.QueryRow(ctx, sql, username).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
	); err != nil {
		return nil, fmt.Errorf("error query row: %w", err)
	}

	return &u, nil
}

func (r *ShortlinkRepository) GetShortlinkUser(
	ctx context.Context,
	link string,
) (*domain.Shortlink, error) {
	sql := `
    SELECT id, user_id, shortlink, original_link, count, created_at
    FROM shortlink.shortlink
    WHERE shortlink = $1
    `

	var sl domain.Shortlink
	if err := r.pool.QueryRow(ctx, sql, link).Scan(
		&sl.ID,
		&sl.UserID,
		&sl.Shortlink,
		&sl.OriginalUrl,
		&sl.Count,
		&sl.CreatedAt,
	); err != nil {
		return &domain.Shortlink{}, fmt.Errorf("failed get short link user: %w", err)
	}

	return &sl, nil
}

func (r *ShortlinkRepository) GetShortLinkByUserID(
	ctx context.Context,
	userID int,
) (*domain.Shortlink, error) {
	sql := `
	SELECT id, user_id, shortlink, original_link, count, created_at
	FROM shortlink.shortlink
	WHERE user_id = $1
	`
  
	var sl domain.Shortlink
	if err := r.pool.QueryRow(ctx, sql, userID).Scan(
		&sl.ID,
		&sl.UserID,
		&sl.Shortlink,
		&sl.OriginalUrl,
		&sl.Count,
		&sl.CreatedAt,
	); err != nil {
		return &domain.Shortlink{}, fmt.Errorf("failed get short link user: %w", err)
	}

	return &sl, nil
}

func (r *ShortlinkRepository) GetShortlinkActivateUser(
	ctx context.Context,
	ip_addr string,
) (*domain.ShortlinkInfo, error) {
	sql := `
   SELECT id, ip_addr, user_agent, shortlink_id, last_activate
   FROM shortlink.shortlink_activate
   WHERE ip_addr = $1
   `

	var s domain.ShortlinkInfo
	if err := r.pool.QueryRow(ctx, sql, ip_addr).Scan(
		&s.ID,
		&s.IPAddr,
		&s.UserAgent,
		&s.ShortlinkID,
		&s.LastAccessAt,
	); err != nil {
		return nil, fmt.Errorf("failed query row: %w", err)
	}

	return &s, nil
}

func (r *ShortlinkRepository) GetShortLinkActivateUsers(
	ctx context.Context,
	ShortlinkID int,
) ([]domain.ShortlinkInfo, error) {
	sql := `
	SELECT id, ip_addr, user_agent, shortlink_id, last_activate
  FROM shortlink.shortlink_activate
  WHERE shortlink_id = $1
	`

	rows, err := r.pool.Query(ctx, sql, ShortlinkID)
	if err != nil {
		return []domain.ShortlinkInfo{}, fmt.Errorf("failed get short link activate user: %w", err)
	}

	shortlink := make([]domain.ShortlinkInfo, 0)

	for rows.Next() {
		var sl domain.ShortlinkInfo

		rows.Scan(
			&sl.ID,
			&sl.IPAddr,
			&sl.UserAgent,
			&sl.ShortlinkID,
			&sl.LastAccessAt,
		)
		defer rows.Close()

		shortlink = append(shortlink, sl)
	}

	return shortlink, nil
}

func (r *ShortlinkRepository) ProfileUser(
	ctx context.Context,
	userID int,
) (*domain.User, error) {
	sql := `
	SELECT id, username, email, register_at
	 FROM shortlink.user
   WHERE id = $1
	`
  
	var user domain.User
	if err := r.pool.QueryRow(ctx, sql, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.RegisterAt,
	); err != nil {
		return &domain.User{}, fmt.Errorf("failed query: %w", err)
	}

	return &user, nil
}