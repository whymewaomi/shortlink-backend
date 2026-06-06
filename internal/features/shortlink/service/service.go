package shortlink_service

import (
	"api/internal/core/domain"
	"context"
	"time"
)

type ShortlinkService struct {
	shortlinkRepository ShortlinkRepository
	shortlinkStorage    ShortlinkStorage
}

type ShortlinkRepository interface {
	CreateShortlink(
		ctx context.Context,
		sl *domain.Shortlink,
	) error
	GetUserByUsername(
		ctx context.Context,
		username string,
	) (*domain.User, error)
	CreateUser(
		ctx context.Context,
		user *domain.User,
	) (int, error)
	GetShortlinkUser(
		ctx context.Context,
		link string,
	) (*domain.Shortlink, error)
	GetShortlinkActivateUser(
		ctx context.Context,
		ip_addr string,
	) (*domain.ShortlinkInfo, error)
	UpdateCountLink(
	ctx context.Context,
	userID int,
	count int,
  ) error
  GetShortLinkActivateUsers(
	ctx context.Context,
	ShortlinkID int,
  ) ([]domain.ShortlinkInfo, error)
  GetShortLinkByUserID(
	ctx context.Context,
	userID int,
  ) (*domain.Shortlink, error)
  CreateShortLinkActivate(
	ctx context.Context,
	sli *domain.ShortlinkInfo,
) error
  ProfileUser(
	ctx context.Context,
	userID int,
) (*domain.User, error)
}

type ShortlinkStorage interface {
	Set(
		ctx context.Context,
		key string,
		value interface{},
		ttl time.Duration,
	) error
	Get(
		ctx context.Context,
		key string,
	) (string, error)
	Del(
		ctx context.Context,
		key string,
	) error
}

func NewShortlinkService(
	shortlinkRepository ShortlinkRepository,
	shortlinkStorage ShortlinkStorage,
) *ShortlinkService {
	return &ShortlinkService{
		shortlinkRepository: shortlinkRepository,
		shortlinkStorage:    shortlinkStorage,
	}
}
