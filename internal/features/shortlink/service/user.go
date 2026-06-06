package shortlink_service

import (
	"api/internal/core/domain"
	core_jwt "api/internal/core/jwt"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *ShortlinkService) HashPassword(
	password string,
) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed generate password: %w", err)
	}

	return string(passwordHash), nil
}

func (s *ShortlinkService) HashRefreshToken() (string, error) {
	refToken := make([]byte, 32)

	if _, err := rand.Read(refToken); err != nil {
		return "", err
	}

	return hex.EncodeToString(refToken), nil
}

func (s *ShortlinkService) SaveRefreshToken(
	ctx context.Context,
	userID int,
	device *domain.Device,
) (*domain.Session, error) {
	refToken, err := s.HashRefreshToken()
	if err != nil {
		return &domain.Session{}, fmt.Errorf("error hash refresh: %w", err)
	}
	keyCache := fmt.Sprintf("session:%s", refToken)

	session := domain.NewSession(userID, refToken, device)
	sessionCache, err := json.Marshal(session)
	if err != nil {
		return &domain.Session{}, fmt.Errorf("failed marshal: %w", err)
	}

	return session, s.shortlinkStorage.Set(ctx, keyCache, string(sessionCache), 30*24*time.Hour)
}

func (s *ShortlinkService) RegisterUser(
	ctx context.Context,
	user *domain.User,
) (*domain.JWT, error) {
	userFromDB, err := s.shortlinkRepository.GetUserByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &domain.JWT{}, fmt.Errorf("failed pgx: %w", err)
	}
	if userFromDB != nil {
		return &domain.JWT{}, fmt.Errorf("username '%s' already exists", user.Username)
	}

	if err := user.Validate(); err != nil {
		return &domain.JWT{}, fmt.Errorf("validate error: %w", err)
	}
	passwordHash, err := s.HashPassword(user.Password)
	if err != nil {
		return &domain.JWT{}, fmt.Errorf("failed hashed password: %w", err)
	}
	user.Password = passwordHash

	userID, err := s.shortlinkRepository.CreateUser(ctx, user)
	if err != nil {
		return &domain.JWT{}, fmt.Errorf("failed create user: %w", err)
	}
	token, err := core_jwt.GenerateJWT(userID)
	if err != nil {
		return &domain.JWT{}, fmt.Errorf("error generate jwt: %w", err)
	}

	return domain.NewJwt(userID, token), nil
}

func (s *ShortlinkService) CheckPassword(
	ctx context.Context,
	user *domain.User,
) (int, error) {
	userFromDB, err := s.shortlinkRepository.GetUserByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed pgx: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password)); err != nil {
		return 0, fmt.Errorf("user or password incorrect: %w", err)
	}

	return userFromDB.ID, nil
}

func (s *ShortlinkService) LoginUser(
	ctx context.Context,
	user *domain.User,
) (*domain.JWT, error) {
	userID, err := s.CheckPassword(ctx, user)
	if err != nil {
		return &domain.JWT{}, fmt.Errorf("error: %w", err)
	}

	token, err := core_jwt.GenerateJWT(userID)
	if err != nil {
		return &domain.JWT{}, fmt.Errorf("error generate jwt: %w", err)
	}

	return domain.NewJwt(userID, token), nil
}

func (s *ShortlinkService) RefreshToken(
	ctx context.Context,
	refreshToken string,
	device *domain.Device,
) (string, error) {
	keyCache := fmt.Sprintf("session:%s", refreshToken)
	sessionCached, err := s.shortlinkStorage.Get(ctx, keyCache)
	if err != nil {
		return "", fmt.Errorf("not found: %w", err)
	}

	var session domain.Session
	if err := json.Unmarshal([]byte(sessionCached), &session); err != nil {
		return "", fmt.Errorf("failed unmarshal: %w", err)
	}

	if device.IpAddr != session.Device.IpAddr || device.UserAgent != session.Device.UserAgent {
		return "", errors.New("error devices")
	}

	return core_jwt.GenerateJWT(session.UserID)
}

func (s *ShortlinkService) LogoutUser(
	ctx context.Context,
	refreshToken string,
) error {
	return s.shortlinkStorage.Del(ctx, fmt.Sprintf("session:%s", refreshToken))
}

func (s *ShortlinkService) ProfileUser(
	ctx context.Context,
	userID int,
) (*domain.User, error) {
	keyCached := fmt.Sprintf("profile:%d", userID)
	profileCached, err := s.shortlinkStorage.Get(ctx, keyCached)
	if err == nil {
		var user domain.User

		if err := json.Unmarshal([]byte(profileCached), &user); err != nil {
			return &domain.User{}, fmt.Errorf("failed unmarshal: %w", err)
		}

		return &user, nil
	}

	user, err := s.shortlinkRepository.ProfileUser(ctx, userID)
	if err != nil {
		return &domain.User{}, fmt.Errorf("error get profile: %w", err)
	}

	if err := s.AddCached(ctx, keyCached, user, 30 * time.Second); err != nil {
		return &domain.User{}, fmt.Errorf("failed cached: %w", err)
	}

	return user, nil
}