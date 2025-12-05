package di

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"adapter/internal/config"
	catalogDomain "adapter/internal/domain/catalog_sync"
	"adapter/internal/domain/permissions" 
	registryHandler "adapter/internal/handlers/registry_sync"
	registryDomain "adapter/internal/domain/registry_sync"
	catalogSyncHandler "adapter/internal/handlers/catalog_sync"
	catalogSyncPorts "adapter/internal/ports/catalog_sync"
	permissionsHandler "adapter/internal/handlers/permissions"
	permissionsPorts "adapter/internal/ports/permissions"
	"adapter/internal/shared/caching"
	db "adapter/internal/shared/database"
	logger "adapter/internal/shared/log"
	// redisClient "adapter/internal/shared/redis"
)

type Container struct {
	Config *config.Config
	DB     *gorm.DB
	// RedisClient         *redis.Client
	CacheService        caching.CacheService
	RegistrySyncHandler *registryHandler.RegistrySyncHandler
	PermissionsHandler  *permissionsHandler.PermissionsHandler
	CatalogSyncHandler  *catalogSyncHandler.CatalogSyncHandler
}

func (c *Container) Shutdown(ctx context.Context) error {
	logger.Info(ctx, "Shutting down container resources...")

	if c.DB != nil {
		if err := db.Close(); err != nil {
			logger.Error(ctx, err, "Failed to close database connection")
		}
	}

	// if c.RedisClient != nil {
	// 	if err := redisClient.Close(); err != nil {
	// 		logger.Error(ctx, err, "Failed to close Redis connection")
	// 	}
	// }

	logger.Info(ctx, "Container shutdown complete")
	return nil
}

func InitContainer() (*Container, error) {
	cfg, err := config.LoadConfig()
	ctx := context.Background()
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to load config: %w", err), "Configuration error")
	}

	database, err := db.Init(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to initialize database: %w", err), "Database initialization error")
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// redisDB, err := redisClient.Init(cfg.RedisURL)
	// if err != nil {
	// 	logger.Fatal(ctx, fmt.Errorf("failed to initialize Redis: %w", err), "Redis initialization error")
	// 	return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	// }

	// Run database migrations using golang-migrate only
	logger.Info(ctx, "Running database migrations...")
	if err := database.AutoMigrate(&catalogSyncPorts.Seller{}, &permissionsPorts.Bap{}, &catalogSyncPorts.SellerCatalogState{}, &permissionsPorts.BapAccessPolicy{}); err != nil {
		logger.Fatal(ctx, err, "Failed to run database migrations")
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}
	logger.Info(ctx, "Database migrations completed successfully")

	// Create instances
	// cacheService := caching.NewRedisCacheService(redisDB)
	
	// Catalog Sync
	sellerRepo := catalogSyncPorts.NewGormRepository(database)
	catalogSyncService := catalogDomain.NewCatalogSyncService(sellerRepo)
	catalogSyncHandler := catalogSyncHandler.NewCatalogSyncHandler(catalogSyncService)

	// Permissions
	permissionsRepo := permissionsPorts.NewGormRepository(database)
	permissionsService := permissions.NewPermissionsService(permissionsRepo)
	permissionsHandler := permissionsHandler.NewPermissionsHandler(permissionsService)

	// ONDC / Registry Sync
	ondcService := registryDomain.NewONDCService(sellerRepo, cfg)
	registrySyncHandler := registryHandler.NewRegistrySyncHandler(ondcService)


	return &Container{
		Config:      cfg,
		DB:          database,
		// RedisClient:         redisDB,
		// CacheService:        cacheService,
		RegistrySyncHandler: registrySyncHandler,
		PermissionsHandler:  permissionsHandler,
		CatalogSyncHandler:  catalogSyncHandler,
	}, err
}
