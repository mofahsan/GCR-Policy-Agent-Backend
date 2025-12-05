package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"adapter/internal/config/di"
	appError "adapter/internal/shared/error"
	logger "adapter/internal/shared/log"
	"adapter/internal/shared/middleware"
)

func main() {
	ctx := context.Background()

	container, err := di.InitContainer()
	if err != nil {
		fmt.Printf("Failed to initialize container: %v\n", err)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: appError.ErrorHandler(),
	})

	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggingMiddleware())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": container.Config.OtelService})
	})

	// routes
	routes := app.Group("/v1")
	routes.Post("/permissions", container.PermissionsHandler.UpdatePermissions)
	routes.Post("/permissions/query", container.PermissionsHandler.QueryPermissions)
	routes.Get("/catalog-sync/sellers/:seller_id", container.CatalogSyncHandler.GetSyncStatus)
	routes.Get("/catalog-sync/pending", container.CatalogSyncHandler.GetPendingCatalogSyncSellers)

	// Internal routes (nested under /v1)
	internal := routes.Group("/internal")
	internal.Post("/registry-sync", container.RegistrySyncHandler.SyncRegistry)

	port := container.Config.Port

	fmt.Printf("\n🚀 Starting %s server\n", container.Config.OtelService)
	fmt.Printf("   Port: %s\n", port)
	fmt.Printf("   Logging: ✅ Enabled (Application-level logging active)\n")
	if container.Config.OtelURL != "" {
		fmt.Printf("   Tracing: ✅ Enabled (OTEL endpoint: %s)\n", container.Config.OtelURL)
	} else {
		fmt.Printf("   Tracing: ⚠️  Disabled (Set OTEL_URL to enable)\n")
	}
	fmt.Printf("   Health Check: http://localhost:%s/health\n", port)
	fmt.Printf("\n")

	logger.Infof(ctx, "Starting %s server on port %s", container.Config.OtelService, port)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Fatal(ctx, err, "Error starting server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Info(ctx, "Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := container.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, err, "Error during container shutdown")
	}

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error(ctx, err, "Server forced to shutdown")
	} else {
		logger.Info(ctx, "Server shutdown complete")
	}
}
