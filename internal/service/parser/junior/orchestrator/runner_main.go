package orchestrator

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// Run запускает ПОЛНЫЙ цикл парсинга
func (s *orchestratorService) Run(ctx context.Context) error {
	logger.Info(ctx, "🚀 Starting full parsing cycle...")
	logger.Info(ctx, "================================================================================")

	// Этап 1: Junior.fhr.ru парсинг
	if err := s.RunJuniorParsing(ctx); err != nil {
		return err
	}

	// Этап 2: Registry.fhr.ru парсинг (для будущего)
	// if err := s.RunRegistryParsing(ctx); err != nil {
	//     return err
	// }

	// Этап 3: Merge данных (для будущего)
	// if err := s.RunMerge(ctx); err != nil {
	//     return err
	// }

	logger.Info(ctx, "================================================================================")
	logger.Info(ctx, "✅ Full parsing cycle completed!")
	return nil
}

// RunRegistryParsing парсит registrynew.fhr.ru (заглушка для будущего)
func (s *orchestratorService) RunRegistryParsing(ctx context.Context) error {
	logger.Info(ctx, "⚠️  RegistryParsing not implemented yet")
	return nil
}

// RunMerge мержит данные из разных источников (заглушка для будущего)
func (s *orchestratorService) RunMerge(ctx context.Context) error {
	logger.Info(ctx, "⚠️  Merge not implemented yet")
	return nil
}
