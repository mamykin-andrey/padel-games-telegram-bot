package handlers

import (
	"main/shared"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpCommandHandler struct {
	bot shared.BotAPI
}

func NewHelpCommandHandler(bot shared.BotAPI) *HelpCommandHandler {
	return &HelpCommandHandler{bot: bot}
}

func (h *HelpCommandHandler) HandleCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Supported commands: /help, /new, /games, /join, /delete")
	h.bot.SendMessage(msg)
}
