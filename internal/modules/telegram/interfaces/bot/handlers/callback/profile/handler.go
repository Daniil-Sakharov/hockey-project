package profile

import (
	"context"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter/keyboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler обрабатывает профиль
type Handler struct {
	presenter      *presenter.Presenter
	keyboard       *keyboard.KeyboardPresenter
	profileService *services.ProfileService
}

// NewHandler создает новый Handler
func NewHandler(
	presenter *presenter.Presenter,
	keyboard *keyboard.KeyboardPresenter,
	profileService *services.ProfileService,
) *Handler {
	return &Handler{
		presenter:      presenter,
		keyboard:       keyboard,
		profileService: profileService,
	}
}

// HandleProfile отображает профиль игрока
func (h *Handler) HandleProfile(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	parts := strings.Split(query.Data, ":")
	if len(parts) != 3 {
		return nil
	}

	playerID := parts[2]

	profile, err := h.profileService.GetProfile(ctx, playerID)
	if err != nil {
		return err
	}

	text, err := h.presenter.RenderProfile(profile)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	msg.ReplyMarkup = h.keyboard.ProfileKeyboard(playerID)

	_, err = bot.Send(msg)
	return err
}
