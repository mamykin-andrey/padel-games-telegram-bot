package main

import (
	"log"
	"main/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler interface {
	HandleCommand(update tgbotapi.Update) bool
}

const botTokenEnvName = "PADEL_BOT_TOKEN"

var bot *tgbotapi.BotAPI // TODO: Wrap with an interface
var registeredHandlers map[string]CommandHandler = make(map[string]CommandHandler)
var newGameHandler *handlers.NewGameCommandHandler

func main() {
	bot = initBot()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	registeredHandlers["join"] = handlers.NewJoinGameCommandHandler(bot)
	registeredHandlers["help"] = handlers.NewHelpCommandHandler(bot)
	registeredHandlers["games"] = handlers.NewActiveGamesCommandHandler(bot)
	newGameHandler = handlers.NewNewGameCommandHandler(bot)
	registeredHandlers["new"] = newGameHandler

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if handleCommand(update) {
			continue
		}
		if newGameHandler.HandleNewGameMessage(update) {
			continue
		}
	}
}

func initBot() *tgbotapi.BotAPI {
	token, ok := getEnvValue(botTokenEnvName)
	if !ok || token == "" {
		log.Panicf("%s is not set", botTokenEnvName)
	}
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	debugLog("Authorized on account:", bot.Self.UserName)
	return bot
}

func handleCommand(update tgbotapi.Update) bool {
	command := removeDigits(update.Message.Command())
	if handler, exists := handlers[command]; exists {
		return handler.HandleCommand(update)
	}
	return false
}

func editMessage(chatId int64, messageId int, newText string) {
	editMessageConfig := tgbotapi.NewEditMessageText(chatId, messageId, newText)
	if _, err := bot.Request(editMessageConfig); err != nil {
		log.Panic(err)
	}
}

func deleteMessage(chatId int64, messageId int) {
	deleteMessageConfig := tgbotapi.NewDeleteMessage(chatId, messageId)
	if _, err := bot.Request(deleteMessageConfig); err != nil {
		log.Panic(err)
	}
}

func sendMessage(msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
