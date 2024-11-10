package main

import (
	"log"
	"main/handlers"
	"main/shared"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler interface {
	HandleCommand(update tgbotapi.Update)
}

const botTokenEnvName = "PADEL_BOT_TOKEN"

var bot shared.BotAPI
var registeredHandlers map[string]CommandHandler = make(map[string]CommandHandler)
var newGameHandler *handlers.NewGameCommandHandler

func main() {
	bot = shared.BindingsBotAPI{
		BindingsBot: initBot(),
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	registeredHandlers["join"] = handlers.NewJoinGameCommandHandler(bot)
	registeredHandlers["help"] = handlers.NewHelpCommandHandler(bot)
	registeredHandlers["games"] = handlers.NewActiveGamesCommandHandler(bot)
	registeredHandlers["delete"] = handlers.NewDeleteGameCommandHandler(bot)
	newGameHandler = handlers.NewNewGameCommandHandler(bot)
	registeredHandlers["new"] = newGameHandler

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		go handleUpdate(update)
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.Message == nil && update.CallbackQuery == nil {
		return
	}
	if handleCommand(update) {
		return
	}
	if newGameHandler.HandleNewGameMessage(update) {
		return
	}
}

func initBot() *tgbotapi.BotAPI {
	token, ok := shared.GetEnvValue(botTokenEnvName)
	if !ok || token == "" {
		log.Panicf("%s is not set", botTokenEnvName)
	}
	var err error
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	shared.DebugLog("Authorized on account:", bot.Self.UserName)
	return bot
}

func handleCommand(update tgbotapi.Update) bool {
	command := parseCommand(update)
	if handler, exists := registeredHandlers[command]; exists {
		handler.HandleCommand(update)
		return true
	}
	return false
}

func parseCommand(update tgbotapi.Update) string {
	var input string
	if update.CallbackQuery != nil {
		input = update.CallbackQuery.Data
	} else {
		input = update.Message.Command()
	}
	return shared.RemoveDigits(input)
}
