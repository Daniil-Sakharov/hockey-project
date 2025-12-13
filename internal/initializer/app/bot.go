package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/closer"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// BotApp представляет приложение Telegram бота
type BotApp struct {
	*initializer.App
}

// NewBotApp создает новый экземпляр BotApp
func NewBotApp(ctx context.Context) (*BotApp, error) {
	baseApp, err := initializer.New(ctx)
	if err != nil {
		return nil, err
	}

	return &BotApp{App: baseApp}, nil
}

// Run запускает бота с обработкой сигналов
func (a *BotApp) Run(ctx context.Context) error {
	logger.Info(ctx, "🤖 Starting Telegram Bot...")

	bot := a.DiContainer.TelegramBot(ctx)

	botCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigChan)

	errChan := make(chan error, 1)
	go func() {
		if err := bot.Start(botCtx); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-sigChan:
		logger.Info(ctx, "📛 Получен сигнал завершения")
		cancel()
		closer.CloseAll(ctx)
		logger.Info(ctx, "✅ Bot stopped gracefully")
		return nil

	case err := <-errChan:
		logger.Error(ctx, "❌ Бот остановлен с ошибкой", zap.Error(err))
		cancel()
		closer.CloseAll(ctx)
		return err
	}
}
