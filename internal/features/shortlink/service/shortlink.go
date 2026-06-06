package shortlink_service

import (
	"api/internal/core/domain"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
)

func (s *ShortlinkService) AddCached(
	ctx context.Context,
	key string,
	val interface{},
	ttl time.Duration,
) error {
	valByte, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return s.shortlinkStorage.Set(ctx, key, valByte, ttl)
}

func (s *ShortlinkService) ValidateLink(link string) error {
	u, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("failed parse link: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("unsupported scheme")
	}

	return nil
}

func (s *ShortlinkService) HashUrl() (string, error) {
	urlByte := make([]byte, 6)

	if _, err := rand.Read(urlByte); err != nil {
		return "", err
	}

	return hex.EncodeToString(urlByte), nil
}

func (s *ShortlinkService) ShortlinkPost(
	ctx context.Context,
	userID int,
	link string,
) (string, error) {
	if link == "" {
		return "", errors.New("bad requests")
	}

	if err := s.ValidateLink(link); err != nil {
		return "", fmt.Errorf("failed validate: %w", err)
	}

	shortLink, err := s.HashUrl()
	if err != nil {
		return "", fmt.Errorf("failed hash link: %w", err)
	}

	return shortLink, s.shortlinkRepository.CreateShortlink(ctx, domain.NewShortLink(userID, shortLink, link))
}

func (s *ShortlinkService) ShortlinkGet(
	ctx context.Context,
	link string,
	device *domain.Device,
) (*domain.Shortlink, error) {
	keyCached := fmt.Sprintf("shortlinl:%s", link)
	shortlinkCached, err := s.shortlinkStorage.Get(ctx, keyCached)
	if err == nil {
		var d domain.Shortlink

		if err := json.Unmarshal([]byte(shortlinkCached), &d); err != nil {
			return &domain.Shortlink{}, fmt.Errorf("error unmarhsal: %w", err)
		}

		return &d, nil
	}
	if !errors.Is(err, redis.Nil) {
		return &domain.Shortlink{}, fmt.Errorf("redis error: %w", err)
	}

	sl, err := s.shortlinkRepository.GetShortlinkUser(ctx, link)
	if err != nil {
		return &domain.Shortlink{}, fmt.Errorf("failed get shortlink: %w", err)
	}

	if err := s.AddCached(ctx, keyCached, sl, 10*time.Hour); err != nil {
		return &domain.Shortlink{}, fmt.Errorf("failed add cahced: %w", err)
	}

	shortlinkActivate, err := s.shortlinkRepository.GetShortlinkActivateUser(ctx, device.IpAddr)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &domain.Shortlink{}, fmt.Errorf("failed get shortlink: %w", err)
	}
	if shortlinkActivate != nil {
		return sl, nil
	}

	shortlink := domain.NewShortLinkInfo(device.IpAddr, device.UserAgent, sl.ID)
	if err := s.shortlinkRepository.CreateShortLinkActivate(ctx, shortlink); err != nil {
		return &domain.Shortlink{}, fmt.Errorf("failed create short link: %w", err)
	}

	return sl, nil
}

func (s *ShortlinkService) ShortLinkActivity(
	ctx context.Context,
	userID int,
) ([]domain.ShortlinkInfo, error) {
	activeCached := fmt.Sprintf("active:%d", userID)
	sliCached, err := s.shortlinkStorage.Get(ctx, activeCached)
	if err == nil {
		var c []domain.ShortlinkInfo

		if err := json.Unmarshal([]byte(sliCached), &c); err != nil {
			return []domain.ShortlinkInfo{}, fmt.Errorf("failed unmarshal: %w", err)
		}

		return c, nil
	}
	if !errors.Is(err, redis.Nil) {
		return []domain.ShortlinkInfo{}, fmt.Errorf("redis error: %w", err)
	}

	sl, err := s.shortlinkRepository.GetShortLinkByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return []domain.ShortlinkInfo{}, fmt.Errorf("failed get shortlink: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.ShortlinkInfo{}, nil
	}

	sli, err := s.shortlinkRepository.GetShortLinkActivateUsers(ctx, sl.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return []domain.ShortlinkInfo{}, fmt.Errorf("failed get: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.ShortlinkInfo{}, nil
	}

	if err := s.AddCached(ctx, activeCached, sli, 10 * time.Minute); err != nil {
		return []domain.ShortlinkInfo{}, fmt.Errorf("failed cached: %w", err)
	}

	return sli, nil
}