package callback

import (
	"context"
	"log"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleSearch обрабатывает callback поиска
func HandleSearch(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 2 {
		return
	}

	command := parts[1]

	switch command {
	case cb.SearchPage:
		handleSearchPagination(r, ctx, bot, query, parts)
	case cb.SearchBackToFilters:
		// Сначала удаляем результаты поиска
		if err := r.SearchHandler().HandleCleanup(ctx, bot, query); err != nil {
			log.Printf("Error cleaning up search results: %v", err)
		}
		// Затем отправляем НОВОЕ сообщение с меню фильтров
		if err := r.FilterHandler().HandleFilterMenuNew(ctx, bot, query); err != nil {
			log.Printf("Error handling filter menu: %v", err)
		}
	case cb.SearchBackToResults:
		// Возврат к результатам поиска из профиля игрока
		if err := r.SearchHandler().HandleBackToResults(ctx, bot, query); err != nil {
			log.Printf("Error handling back to results: %v", err)
		}
	default:
		log.Printf("Unknown search command: %s", command)
	}
}

// handleSearchPagination обрабатывает пагинацию результатов поиска
func handleSearchPagination(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 3 {
		return
	}

	direction := parts[2]

	switch direction {
	case cb.PageNext:
		if err := r.SearchHandler().HandlePageNext(ctx, bot, query); err != nil {
			log.Printf("Error handling page next: %v", err)
		}
	case cb.PagePrev:
		if err := r.SearchHandler().HandlePagePrev(ctx, bot, query); err != nil {
			log.Printf("Error handling page prev: %v", err)
		}
	default:
		log.Printf("Unknown pagination direction: %s", direction)
	}
}
