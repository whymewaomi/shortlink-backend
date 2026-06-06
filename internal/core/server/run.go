package core_server

import (
	core_config "api/internal/core/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

type App struct {
	cfg *core_config.Config

	app *fiber.App
}

func NewApp(
	cfg *core_config.Config,
) *App {
	return &App{
		cfg: cfg,
		app: fiber.New(),
	}
}

func (a *App) GetApp() *fiber.App {
	return a.app
}

func (a *App) Run(ctx context.Context) error {
	ch := make(chan error, 1)

	go func(){
		defer close(ch)
    
		host := fmt.Sprintf(":%s", a.cfg.Host)
		if err := a.app.Listen(host); err != nil {
			ch <- err
			return 
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("failed to start HTTP server: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
		defer cancel()

		log.Printf("Shutdown HTTP server...")

		if err := a.app.ShutdownWithContext(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server error: %w", err)
		}
	}

	return nil
}

