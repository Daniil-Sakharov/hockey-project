package search

import (
	"context"
	"log"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/message"
	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// showSearchResults показывает результаты поиска
func (h *Handler) showSearchResults(ctx context.Context, botAPI *tgbotapi.BotAPI, chatID, userID int64) error {
	state := h.stateManager.GetState(userID)

	// Удаляем старые сообщения с результатами (если есть)
	h.deleteOldResults(botAPI, chatID, state)

	// Выполняем поиск (5 игроков на страницу)
	result, err := h.searchService.Search(ctx, state.Filters, state.CurrentPage, 5)
	if err != nil {
		log.Printf("Error searching players: %v", err)
		return err
	}

	// Если результатов нет - показываем сообщение
	if result.TotalCount == 0 {
		msg := tgbotapi.NewMessage(chatID, "🔍 Игроки не найдены.\n\nПопробуйте изменить фильтры.")
		sent, _ := botAPI.Send(msg)
		state.SearchResultMessageIDs = []int{sent.MessageID}
		return nil
	}

	// Отправляем карточки игроков (5 сообщений)
	messageIDs := []int{}
	for _, player := range result.Players {
		emoji := presenter.GetTeamEmoji(player.TeamName)
		cardText := message.RenderPlayerCard(player, emoji)

		msg := tgbotapi.NewMessage(chatID, cardText)
		markup := h.keyboardPresenter.PlayerProfile(player.Player.ID)
		msg.ReplyMarkup = markup

		sent, err := botAPI.Send(msg)
		if err != nil {
			log.Printf("Error sending player card: %v", err)
			continue
		}
		messageIDs = append(messageIDs, sent.MessageID)
	}

	// Отправляем сообщение с пагинацией (6-е сообщение)
	paginationText := h.msgPresenter.RenderSearchPagination(result)
	paginationMsg := tgbotapi.NewMessage(chatID, paginationText)
	paginationMarkup := h.keyboardPresenter.SearchPagination(result.CurrentPage, result.TotalPages)
	paginationMsg.ReplyMarkup = paginationMarkup

	sent, err := botAPI.Send(paginationMsg)
	if err != nil {
		log.Printf("Error sending pagination: %v", err)
		return err
	}
	messageIDs = append(messageIDs, sent.MessageID)

	// Сохраняем все message_id для будущего удаления
	state.SearchResultMessageIDs = messageIDs

	return nil
}

// deleteOldResults удаляет старые сообщения результатов поиска параллельно
func (h *Handler) deleteOldResults(botAPI *tgbotapi.BotAPI, chatID int64, state *domainBot.UserState) {
	if len(state.SearchResultMessageIDs) == 0 {
		return
	}

	// Используем goroutines для параллельного удаления
	var wg sync.WaitGroup
	for _, msgID := range state.SearchResultMessageIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, id)
			if _, err := botAPI.Request(deleteMsg); err != nil {
				log.Printf("Error deleting old result message %d: %v", id, err)
			}
		}(msgID)
	}

	// Ждем завершения всех удалений
	wg.Wait()

	// Очищаем список
	state.SearchResultMessageIDs = []int{}
}
