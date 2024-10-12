package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpCommandHandler struct {
	bot *tgbotapi.BotAPI
}

func NewHelpCommandHandler(bot *tgbotapi.BotAPI) *HelpCommandHandler {
	return &HelpCommandHandler{bot: bot}
}

func (h *HelpCommandHandler) HandleCommand(update tgbotapi.Update) bool {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Supported commands: /help, /new, /games")
	h.bot.Send(msg)
	return true
}
