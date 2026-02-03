package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/application/services"
	router "github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/interfaces/http"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/interfaces/http/handlers"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/interfaces/http/middleware"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	logConfig := getLoggerConfig()
	if err := logger.Init(getEnv("LOG_LEVEL", "info"), true, logConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info(ctx, "Starting Hockey API...")

	// DI Container
	container := di.NewContainer()
	defer func() { _ = container.Close() }()

	// Database
	db, err := container.DB(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", zap.Error(err))
	}
	logger.Info(ctx, "Connected to database")

	// Auth config
	authConfig := services.AuthConfig{
		JWTSecret:            getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour, // 7 days
	}

	// Services
	statsService := services.NewStatsService(db)
	rankingService := services.NewRankingService(db)
	authService := services.NewAuthService(db, authConfig)
	exploreService := services.NewExploreService(db)
	explorePlayersService := services.NewExplorePlayersService(db)
	exploreMatchesService := services.NewExploreMatchesService(db)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	statsHandler := handlers.NewStatsHandler(statsService)
	rankingHandler := handlers.NewRankingHandler(rankingService)
	authHandler := handlers.NewAuthHandler(authService)
	exploreHandler := handlers.NewExploreHandler(exploreService, exploreMatchesService)
	explorePlayersHandler := handlers.NewExplorePlayersHandler(explorePlayersService)
	exploreMatchesHandler := handlers.NewExploreMatchesHandler(exploreMatchesService)
	imageProxyHandler := handlers.NewImageProxyHandler()

	// Router
	allowedOrigins := []string{"*"} // TODO: configure from env
	apiRouter := router.NewRouter(
		healthHandler,
		statsHandler,
		rankingHandler,
		authHandler,
		exploreHandler,
		explorePlayersHandler,
		exploreMatchesHandler,
		imageProxyHandler,
		authMiddleware,
		allowedOrigins,
	)
	handler := apiRouter.Setup()

	// HTTP Server
	port := getEnv("API_PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Info(ctx, "Received shutdown signal")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error(ctx, "Server shutdown error: "+err.Error())
		}
		cancel()
	}()

	logger.Info(ctx, "API server starting on :"+port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal(ctx, "Server failed", zap.Error(err))
	}

	logger.Info(ctx, "Server stopped")
}

func getLoggerConfig() *logger.LoggerConfig {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil
	}

	return &logger.LoggerConfig{
		ServiceName:  getEnv("OTEL_SERVICE_NAME", "hockey-api"),
		Environment:  getEnv("APP_ENV", "development"),
		OTLPEndpoint: endpoint,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
