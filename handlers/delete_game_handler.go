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

func (h *DeleteGameCommandHandler) HandleCommand(update tgbotapi.Update) {
	command := update.CallbackQuery.Data
	chatId := update.CallbackQuery.Message.Chat.ID
	gameId, _ := strconv.Atoi(command[4:])
	userId := update.CallbackQuery.From.ID
	for i, g := range shared.State.Games() {
		if g.Id == gameId && g.CreatorId == userId {
			shared.State.Remove(i)
			h.bot.SendMessage(tgbotapi.NewMessage(chatId, "The game has been deleted"))
			NewActiveGamesCommandHandler(h.bot).ShowAllGames(chatId)
			return
		}
	}
	h.bot.SendMessage(tgbotapi.NewMessage(chatId, "Only the creator can delete a game"))
}
