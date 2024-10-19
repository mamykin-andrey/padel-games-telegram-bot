package shared

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI interface {
	EditMessage(chatId int64, messageId int, newText string)
	DeleteMessage(chatId int64, messageId int)
	SendMessage(msg tgbotapi.MessageConfig)
	ID() int64
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
}

type BindingsBotAPI struct {
	BindingsBot *tgbotapi.BotAPI
}

func (bot BindingsBotAPI) EditMessage(chatId int64, messageId int, newText string) {
	editMessageConfig := tgbotapi.NewEditMessageText(chatId, messageId, newText)
	if _, err := bot.BindingsBot.Request(editMessageConfig); err != nil {
		log.Panic(err)
	}
}

func (bot BindingsBotAPI) DeleteMessage(chatId int64, messageId int) {
	deleteMessageConfig := tgbotapi.NewDeleteMessage(chatId, messageId)
	if _, err := bot.BindingsBot.Request(deleteMessageConfig); err != nil {
		log.Panic(err)
	}
}

func (bot BindingsBotAPI) SendMessage(msg tgbotapi.MessageConfig) {
	if _, err := bot.BindingsBot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (bot BindingsBotAPI) ID() int64 {
	return bot.BindingsBot.Self.ID
}

func (bot BindingsBotAPI) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return bot.BindingsBot.GetUpdatesChan(config)
}
