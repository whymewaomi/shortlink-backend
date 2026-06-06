package main

import (
	core_config "api/internal/core/config"
	core_pgx "api/internal/core/repository/pgx"
	core_redis "api/internal/core/repository/redis"
	core_server "api/internal/core/server"
	shortlink_reposiotry "api/internal/features/shortlink/repository"
	shortlink_service "api/internal/features/shortlink/service"
	shortlink_transport_http "api/internal/features/shortlink/transport/http"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	cfg := core_config.Load()

	if err := cfg.Validate(); err != nil {
		log.Printf("failed config load: %v", err)
	}

	pool, err := core_pgx.NewConnect(ctx, cfg)
	if err != nil {
		log.Fatalf("error connect pgx: %v", err)
	}
	rdb := core_redis.NewRedisClient(cfg)

	server := core_server.NewApp(cfg)
	app := server.GetApp()

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	shortlinkRepository := shortlink_reposiotry.NewShortLinkRepository(pool)
	shortlinkService := shortlink_service.NewShortlinkService(shortlinkRepository, rdb)
	shortlinkTransport := shortlink_transport_http.NewShortlinkHTTPHandler(shortlinkService, app)
	shortlinkTransport.RegisterRouter()

	if err := server.Run(ctx); err != nil {
		log.Fatalf("error started HTTP server: %v", err)
	}
}
