package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	fhspbRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/retry"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/postgres"
	"go.uber.org/zap"
)

func main() {
	// Загружаем конфигурацию
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	cfg := config.Get()

	// Инициализируем логгер
	if err := logger.Init(cfg.Logger); err != nil {
		log.Fatal("Failed to init logger:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info(ctx, "🔄 Starting FHSPB retry processor...")

	// Подключаемся к PostgreSQL
	db, err := postgres.Connect(cfg.Postgres)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer db.Close()

	logger.Info(ctx, "✅ PostgreSQL connected")

	// Инициализируем клиент и репозитории
	client := fhspb.NewClient()
	playerRepo := fhspbRepo.NewPlayerRepository(db)
	playerTeamRepo := fhspbRepo.PlayerTeamRepository{}

	// Инициализируем retry manager
	retryManager := retry.NewManager(db, cfg.FHSPB.RetryMaxAttempts(), cfg.FHSPB.RetryDelay())

	// Обрабатываем retry задачи
	processor := &RetryProcessor{
		client:         client,
		playerRepo:     playerRepo,
		playerTeamRepo: &playerTeamRepo,
		retryManager:   retryManager,
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info(ctx, "🛑 Shutdown signal received")
		cancel()
	}()

	// Запускаем обработку retry задач
	if err := processor.ProcessRetryJobs(ctx); err != nil {
		logger.Fatal(ctx, "Retry processing failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ FHSPB retry processor completed")
}

type RetryProcessor struct {
	client         *fhspb.Client
	playerRepo     *fhspbRepo.PlayerRepository
	playerTeamRepo *fhspbRepo.PlayerTeamRepository
	retryManager   *retry.Manager
}

func (p *RetryProcessor) ProcessRetryJobs(ctx context.Context) error {
	jobs, err := p.retryManager.GetJobsForRetry(ctx, "fhspb", 100)
	if err != nil {
		return err
	}

	logger.Info(ctx, "📋 Found retry jobs", zap.Int("count", len(jobs)))

	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		logger.Info(ctx, "🔄 Processing retry job", 
			zap.String("type", string(job.JobType)),
			zap.String("external_id", job.ExternalID),
			zap.Int("retry_count", job.RetryCount))

		success := p.processJob(ctx, job)
		
		if err := p.retryManager.MarkJobRetried(ctx, job.ID, success, nil); err != nil {
			logger.Error(ctx, "Failed to mark job as retried", zap.Error(err))
		}
	}

	return nil
}

func (p *RetryProcessor) processJob(ctx context.Context, job retry.FailedJob) bool {
	switch job.JobType {
	case retry.JobTypePlayer:
		return p.processPlayerJob(ctx, job)
	default:
		logger.Warn(ctx, "Unknown job type", zap.String("type", string(job.JobType)))
		return false
	}
}

func (p *RetryProcessor) processPlayerJob(ctx context.Context, job retry.FailedJob) bool {
	// Здесь можно добавить логику повторной обработки игрока
	// Пока просто логируем
	logger.Info(ctx, "🔄 Retrying player job", zap.String("player_id", job.ExternalID))
	return true
}
