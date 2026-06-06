package shortlink_transport_http

import (
	"api/internal/core/domain"
	core_dto "api/internal/core/dto"
	"api/internal/core/validator"

	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

const userId = "user_id"

func (h *ShortlinkHTTPHandler) RegisterUser(c fiber.Ctx) error {
		var r core_dto.RegisterUserDto
		if err := c.Bind().Body(&r); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := validator.Valide(&r); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if r.Password != r.PasswordConfirm {
			return c.Status(400).JSON(fiber.Map{
				"error": "passwords failed",
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		user := domain.NewUser(r.Username, r.Email, r.Password)
		jwt, err := h.shortlinkService.RegisterUser(ctx, user)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		device := domain.NewDevice(c.IP(), c.UserAgent())
		session, err := h.shortlinkService.SaveRefreshToken(ctx, jwt.UserID, device)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    session.RefreshToken,
			Path:     "/",
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			MaxAge:   60 * 60 * 24 * 30,
		})

		return c.Status(201).JSON(core_dto.UserDtoResponse{
			AccessToken:  jwt.Token,
			RefreshToken: session.RefreshToken,
		})
	}

func (h *ShortlinkHTTPHandler) LoginUser(c fiber.Ctx) error {
		var u core_dto.LoginUserDto
		if err := c.Bind().Body(&u); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := validator.Valide(&u); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		user := domain.NewLoginUser(u.Username, u.Password)
		jwt, err := h.shortlinkService.LoginUser(ctx, user)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		device := domain.NewDevice(c.IP(), c.UserAgent())
		session, err := h.shortlinkService.SaveRefreshToken(ctx, jwt.UserID, device)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Bad request",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    session.RefreshToken,
			Path:     "/",
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			MaxAge:   60 * 60 * 24 * 30,
		})

		return c.Status(200).JSON(core_dto.UserDtoResponse{
			AccessToken:  jwt.Token,
			RefreshToken: session.RefreshToken,
		})
	}

func (h *ShortlinkHTTPHandler) RefreshTokenHandler(c fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			return c.Status(404).JSON(fiber.Map{
				"error": "token not found",
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		device := domain.NewDevice(c.IP(), c.UserAgent())
		jwt, err := h.shortlinkService.RefreshToken(ctx, refreshToken, device)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Bad request",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"access_token": jwt,
		})
	}

func (s *ShortlinkHTTPHandler) LogoutHandler(c fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "token not found",
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		if err := s.shortlinkService.LogoutUser(ctx, refreshToken); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Bad request",
			})
		}

		c.ClearCookie("refresh_token")

		return c.Status(200).JSON(fiber.Map{
			"status": "OK",
		})
	}

func (h *ShortlinkHTTPHandler) ProfileUserHandler(c fiber.Ctx) error {
		userID, ok := c.Locals(userId).(int)
		if !ok {
      return c.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
			})
    }

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		user, err := h.shortlinkService.ProfileUser(ctx, userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(user)
	}