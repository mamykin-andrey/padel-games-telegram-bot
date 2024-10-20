package shared

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI interface {
	EditMessage(chatId int64, messageId int, newText string)
	DeleteMessage(chatId int64, messageId int)
	SendMessage(msg tgbotapi.MessageConfig)
	EditMessageTextAndMarkup(chatId int64, messageId int, newText string, replyMarkup tgbotapi.InlineKeyboardMarkup)
	ID() int64
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	EditMessageTextAndRemoveMarkup(chatId int64, messageId int, newText string)
}

type BindingsBotAPI struct {
	BindingsBot *tgbotapi.BotAPI
}

func (bot BindingsBotAPI) EditMessage(chatId int64, messageId int, newText string) {
	config := tgbotapi.NewEditMessageText(chatId, messageId, newText)
	if _, err := bot.BindingsBot.Request(config); err != nil {
		log.Panic(err)
	}
}

func (bot BindingsBotAPI) EditMessageTextAndRemoveMarkup(chatId int64, messageId int, newText string) {
	config := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatId,
			MessageID:   messageId,
			ReplyMarkup: nil,
		},
		Text: newText,
	}
	if _, err := bot.BindingsBot.Request(config); err != nil {
		log.Panic(err)
	}
}

func (bot BindingsBotAPI) EditMessageTextAndMarkup(chatId int64, messageId int, newText string, replyMarkup tgbotapi.InlineKeyboardMarkup) {
	config := tgbotapi.NewEditMessageTextAndMarkup(chatId, messageId, newText, replyMarkup)
	if _, err := bot.BindingsBot.Request(config); err != nil {
		log.Panic(err)
	}
}

func (bot BindingsBotAPI) DeleteMessage(chatId int64, messageId int) {
	config := tgbotapi.NewDeleteMessage(chatId, messageId)
	if _, err := bot.BindingsBot.Request(config); err != nil {
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
