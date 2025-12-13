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

func (h *Handler) showSearchResults(ctx context.Context, botAPI *tgbotapi.BotAPI, chatID, userID int64) error {
	state := h.stateManager.GetState(userID)

	h.deleteOldResults(botAPI, chatID, state)

	result, err := h.searchService.Search(ctx, state.Filters, state.CurrentPage, 5)
	if err != nil {
		log.Printf("Error searching players: %v", err)
		return err
	}

	if result.TotalCount == 0 {
		msg := tgbotapi.NewMessage(chatID, "🔍 Игроки не найдены.\n\nПопробуйте изменить фильтры.")
		sent, _ := botAPI.Send(msg)
		state.SearchResultMessageIDs = []int{sent.MessageID}
		return nil
	}

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

	state.SearchResultMessageIDs = messageIDs
	return nil
}

func (h *Handler) deleteOldResults(botAPI *tgbotapi.BotAPI, chatID int64, state *domainBot.UserState) {
	if len(state.SearchResultMessageIDs) == 0 {
		return
	}

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

	wg.Wait()
	state.SearchResultMessageIDs = []int{}
}
