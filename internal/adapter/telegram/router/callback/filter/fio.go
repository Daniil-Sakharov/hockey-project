package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleFio обрабатывает фильтр ФИО
func HandleFio(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string, logger *zap.Logger) {
	if len(parts) < 3 {
		return
	}

	subCommand := parts[2]

	switch subCommand {
	case cb.FioSelect:
		// Открыть панель ФИО
		if err := r.FilterHandler().HandleFioSelect(ctx, bot, query); err != nil {
			logger.Error("Error handling FIO select", zap.Error(err))
		}
	case cb.FioLastName, cb.FioFirstName, cb.FioPatronymic:
		// Запустить ввод поля
		if err := r.FilterHandler().HandleFioField(ctx, bot, query, subCommand); err != nil {
			logger.Error("Error handling FIO field",
				zap.Error(err),
				zap.String("field", subCommand))
		}
	case cb.FioClearLast, cb.FioClearFirst, cb.FioClearPatr:
		// Очистить поле
		if err := r.FilterHandler().HandleFioClear(ctx, bot, query, subCommand); err != nil {
			logger.Error("Error handling FIO clear",
				zap.Error(err),
				zap.String("field", subCommand))
		}
	case cb.FioApply:
		// Применить изменения
		if err := r.FilterHandler().HandleFioApply(ctx, bot, query); err != nil {
			logger.Error("Error handling FIO apply", zap.Error(err))
		}
	case cb.FioBack:
		// Отменить и вернуться
		if err := r.FilterHandler().HandleFioBack(ctx, bot, query); err != nil {
			logger.Error("Error handling FIO back", zap.Error(err))
		}
	default:
		logger.Warn("Unknown FIO command", zap.String("command", subCommand))
	}
}
