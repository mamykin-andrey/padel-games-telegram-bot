package main

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewGameState int

const (
	NotStarted NewGameState = iota
	Started
	AwaitDate
	AwaitTime
	AwaitDuration
	AwaitPlace
	AwaitPlayers
	AwaitLevel
)

const botTokenEnvName = "PADEL_BOT_TOKEN"

var userGameStates = make(map[int64]NewGameState)
var games []Game
var currentGameId = 0
var bot *tgbotapi.BotAPI // TODO: Wrap with an interface

func main() {
	initBot()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil || (!update.Message.IsCommand() && !isUserCreatingGame(update.Message.From.ID)) {
			continue
		}
		if handleCommand(update) {
			continue
		}
		if handleNewGameMessage(update) {
			continue
		}
	}
}

func isUserCreatingGame(userId int64) bool {
	_, exists := userGameStates[userId]
	return exists
}

func initBot() {
	token, ok := getEnvValue(botTokenEnvName)
	if !ok || token == "" {
		log.Panic("PADEL_BOT_TOKEN is not set")
	}
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	debugLog("Authorized on account:", bot.Self.UserName)
}

func handleNewGameMessage(update tgbotapi.Update) bool {
	if !isReplyToTheBot(update) || !isUserCreatingGame(update.Message.From.ID) || update.Message.IsCommand() {
		return false
	}
	return transitionGameState(update)
}

func isReplyToTheBot(update tgbotapi.Update) bool {
	if update.Message.ReplyToMessage == nil {
		return false
	}
	return bot.Self.ID == update.Message.ReplyToMessage.From.ID
}

func transitionGameState(update tgbotapi.Update) bool {
	userId := update.Message.From.ID
	input := update.Message.Text
	gameState := userGameStates[userId]
	userMessageId := update.Message.MessageID
	chatId := update.Message.Chat.ID
	if gameState == NotStarted {
		return false
	}
	if gameState == Started {
		games = append(games, Game{Id: currentGameId})
		currentGameId++
	}
	game := &games[len(games)-1]
	switch gameState {
	case Started:
		userGameStates[userId] = AwaitDate
		msg := tgbotapi.NewMessage(chatId, "Please enter the game date")
		sendMessage(msg)
		return true
	case AwaitDate:
		userGameStates[userId] = AwaitTime
		game.Date = input
		editMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game time")
		deleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitTime:
		userGameStates[userId] = AwaitDuration
		game.Time = input
		editMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game duration")
		deleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitDuration:
		userGameStates[userId] = AwaitPlace
		game.Duration = input
		editMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game place")
		deleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitPlace:
		userGameStates[userId] = AwaitPlayers
		game.Place = input
		editMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter how many spots you have")
		deleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitPlayers:
		userGameStates[userId] = AwaitLevel
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > 3 {
			sendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter a correct number"))
			return true
		}
		game.NumberOfSpots = num
		editMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game level")
		deleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitLevel:
		userGameStates[userId] = NotStarted
		game.Level = input
		game.IsPublished = true
		game.Players = append(game.Players, fmt.Sprint("@", update.Message.From.UserName))
		deleteMessage(update.Message.Chat.ID, userMessageId)
		deleteMessage(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
		handleActiveGames(update)
		return true
	}
	return false
}

func editMessage(chatId int64, messageId int, newText string) {
	editMessageConfig := tgbotapi.NewEditMessageText(chatId, messageId, newText)
	bot.Request(editMessageConfig)
}

func deleteMessage(chatId int64, messageId int) {
	deleteMessageConfig := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Request(deleteMessageConfig)
}

func handleCommand(update tgbotapi.Update) bool {
	if strings.HasPrefix(update.Message.Command(), "join") {
		handleJoinGame(update)
		return true
	}
	switch update.Message.Command() {
	case "help":
		handleHelp(update)
		return true
	case "new":
		handleNewGame(update)
		return true
	case "games":
		handleActiveGames(update)
		return true
	}
	return false
}

func handleNewGame(update tgbotapi.Update) {
	userGameStates[update.Message.From.ID] = Started
	transitionGameState(update)
}

func handleHelp(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Supported commands: /help to show this message, /new to create a new game, /games to show all active games"
	sendMessage(msg)
}

func handleJoinGame(update tgbotapi.Update) {
	if len(games) == 0 {
		return
	}
	gameId, _ := strconv.Atoi(update.Message.Command()[4:])
	game := &games[slices.IndexFunc(games, func(g Game) bool { return g.Id == gameId })]
	userName := fmt.Sprint("@", update.Message.From.UserName)
	if isPlayerInGame(*game, userName) || game.isFull() {
		return
	}
	game.Players = append(game.Players, userName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Joined the game")
	sendMessage(msg)
}

func isPlayerInGame(game Game, userName string) bool {
	for _, p := range game.Players {
		if p == userName {
			return true
		}
	}
	return false
}

func handleActiveGames(update tgbotapi.Update) {
	activeGames := make([]Game, 0)
	for _, g := range games {
		if g.IsPublished {
			activeGames = append(activeGames, g)
		}
	}
	if len(activeGames) == 0 {
		sendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "No active games"))
	} else {
		for _, g := range activeGames {
			gamePlayers := strings.Join(g.Players[:], ", ")
			gameStr := fmt.Sprint(
				"📅 Date: ", g.Date,
				"\n⏰ Time: ", g.Time,
				"\n⏲️ Duration: ", g.Duration,
				"\n📊 Level: ", g.Level,
				"\n📍 Location: ", g.Place,
				"\n🏋🏻‍♂️ Players: ", 4-g.NumberOfSpots, " + ", gamePlayers,
				"\nJoin the game: /join", g.Id,
			)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, gameStr)
			sendMessage(msg)
		}
	}
}

func sendMessage(msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
