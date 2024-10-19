package handlers

import (
	"main/shared"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeleteGameCommandHandler struct {
	bot shared.BotAPI
}

func NewDeleteGameCommandHandler(bot shared.BotAPI) *DeleteGameCommandHandler {
	return &DeleteGameCommandHandler{bot: bot}
}

func (h *DeleteGameCommandHandler) HandleCommand(update tgbotapi.Update) bool {
	gameId, _ := strconv.Atoi(update.Message.Command()[4:])
	userId := update.Message.From.ID
	for i, g := range shared.Games {
		if g.Id == gameId && g.CreatorId == userId {
			shared.Games = append(shared.Games[:i], shared.Games[i+1:]...)
		}
	}
	h.bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "The game has been deleted"))
	return NewActiveGamesCommandHandler(h.bot).HandleCommand(update)
}
