package shortlink_transport_http

import (
	"api/internal/core/domain"
	core_middleware "api/internal/core/middleware"
	"context"

	"github.com/gofiber/fiber/v3"
)

type ShortlinkHTTPHandler struct {
	shortlinkService ShortinkService

	app *fiber.App
}

type ShortinkService interface {
	ShortlinkPost(
		ctx context.Context,
		userID int,
		link string,
	) (string, error)
	RegisterUser(
		ctx context.Context,
		user *domain.User,
	) (*domain.JWT, error)
	LoginUser(
		ctx context.Context,
		user *domain.User,
	) (*domain.JWT, error)
	SaveRefreshToken(
		ctx context.Context,
		userID int,
		device *domain.Device,
	) (*domain.Session, error)
	RefreshToken(
		ctx context.Context,
		refreshToken string,
		device *domain.Device,
	) (string, error)
	LogoutUser(
		ctx context.Context,
		refreshToken string,
	) error
	ShortlinkGet(
		ctx context.Context,
		link string,
		device *domain.Device,
	) (*domain.Shortlink, error)
	ShortLinkActivity(
	ctx context.Context,
	userID int,
) ([]domain.ShortlinkInfo, error)
  ProfileUser(
	ctx context.Context,
	userID int,
) (*domain.User, error)
}

func NewShortlinkHTTPHandler(
	shortlinkService ShortinkService,
	app *fiber.App,
) *ShortlinkHTTPHandler {
	return &ShortlinkHTTPHandler{
		shortlinkService: shortlinkService,
		app:           app,
	}
}

func (h *ShortlinkHTTPHandler) RegisterRouter() {
	api := h.app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", h.RegisterUser)
	auth.Post("/login", h.LoginUser)
	auth.Post("/refresh", h.RefreshTokenHandler)
	auth.Post("/logout", core_middleware.JWTCheck(), h.LogoutHandler)

  user := api.Group("/user")
	user.Get("/profile", core_middleware.JWTCheck(), h.ProfileUserHandler)

	shortlink := api.Group("/link")
	shortlink.Get("/shortlink/:link", h.ShortLinkGetHandler)
	shortlink.Get("/activate", core_middleware.JWTCheck(), h.ShortLinkActiviteHandler)
	shortlink.Post("/shortlink", core_middleware.JWTCheck(), h.ShortLinkHandler)
}
