package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var currentGameId = 0

type Game struct {
	Id            int
	Date          string
	Time          string
	Duration      string
	Place         string
	Level         string
	Players       []string
	NumberOfSpots int
	IsPublished   bool
}

func (g Game) String() string {
	return fmt.Sprint("Id: ", g.Id, ", date: ", g.Date, ", time: ", g.Time, ", duration: ", g.Duration, ", place: ", g.Place, ", level: ", g.Level, ", players: ", g.Players)
}

var games []Game

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

var userGameStates = make(map[int64]NewGameState)

func main() {
	bot := initBot()
	bot.Debug = true

	debugLog(fmt.Sprint("Authorized on account: ", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil || (!update.Message.IsCommand() && !isUserCreatingGame(update.Message.From.ID)) {
			continue
		}
		if handleCommand(bot, update) {
			continue
		}
		if handleNewGameMessage(bot, update) {
			continue
		}
	}
}

func debugLog(message string) {
	// TODO: Add varargs for the client convenience
	log.Print(message)
}

func isUserCreatingGame(userId int64) bool {
	_, exists := userGameStates[userId]
	return exists
}

func initBot() *tgbotapi.BotAPI {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	token := os.Getenv("PADEL_BOT_TOKEN")
	if token == "" {
		log.Fatal("PADEL_BOT_TOKEN is not set")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return bot
}

func handleNewGameMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	if !isReplyToTheBot(bot, update) || !isUserCreatingGame(update.Message.From.ID) || update.Message.IsCommand() {
		return false
	}
	return transitionGameState(bot, update)
}

func isReplyToTheBot(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	if update.Message.ReplyToMessage == nil {
		return false
	}
	return bot.Self.ID == update.Message.ReplyToMessage.From.ID
}

func transitionGameState(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
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
		sendMessage(bot, msg)
		return true
	case AwaitDate:
		userGameStates[userId] = AwaitTime
		game.Date = input
		editMessage(bot, chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game time")
		deleteMessage(bot, update.Message.Chat.ID, userMessageId)
		return true
	case AwaitTime:
		userGameStates[userId] = AwaitDuration
		game.Time = input
		editMessage(bot, chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game duration")
		deleteMessage(bot, update.Message.Chat.ID, userMessageId)
		return true
	case AwaitDuration:
		userGameStates[userId] = AwaitPlace
		game.Duration = input
		editMessage(bot, chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game place")
		deleteMessage(bot, update.Message.Chat.ID, userMessageId)
		return true
	case AwaitPlace:
		userGameStates[userId] = AwaitPlayers
		game.Place = input
		editMessage(bot, chatId, update.Message.ReplyToMessage.MessageID, "Please enter how many spots you have")
		deleteMessage(bot, update.Message.Chat.ID, userMessageId)
		return true
	case AwaitPlayers:
		userGameStates[userId] = AwaitLevel
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > 3 {
			// TODO: add error handling
			return true
		}
		game.NumberOfSpots = num
		editMessage(bot, chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game level")
		deleteMessage(bot, update.Message.Chat.ID, userMessageId)
		return true
	case AwaitLevel:
		userGameStates[userId] = NotStarted
		game.Level = input
		game.IsPublished = true
		game.Players = append(game.Players, fmt.Sprint("@", update.Message.From.UserName))
		deleteMessage(bot, update.Message.Chat.ID, userMessageId)
		deleteMessage(bot, update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
		handleActiveGames(bot, update)
		return true
	}
	return false
}

func editMessage(bot *tgbotapi.BotAPI, chatId int64, messageId int, newText string) {
	editMessageConfig := tgbotapi.NewEditMessageText(chatId, messageId, newText)
	bot.Request(editMessageConfig)
}

func deleteMessage(bot *tgbotapi.BotAPI, chatId int64, messageId int) {
	deleteMessageConfig := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Request(deleteMessageConfig)
}

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	if strings.HasPrefix(update.Message.Command(), "join") {
		handleJoinGame(bot, update)
		return true
	}
	switch update.Message.Command() {
	case "help":
		handleHelp(bot, update)
		return true
	case "new":
		handleNewGame(bot, update)
		return true
	case "games":
		handleActiveGames(bot, update)
		return true
	}
	return false
}

func handleNewGame(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userGameStates[update.Message.From.ID] = Started
	transitionGameState(bot, update)
}

func handleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Supported commands: /help to show this message, /new to create a new game, /games to show all active games"
	sendMessage(bot, msg)
}

func handleJoinGame(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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
	sendMessage(bot, msg)
}

func (g Game) isFull() bool {
	return len(g.Players) == g.NumberOfSpots
}

func isPlayerInGame(game Game, userName string) bool {
	for _, p := range game.Players {
		if p == userName {
			return true
		}
	}
	return false
}

func handleActiveGames(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	activeGames := make([]Game, 0)
	for _, g := range games {
		if g.IsPublished {
			activeGames = append(activeGames, g)
		}
	}
	if len(activeGames) == 0 {
		sendMessage(bot, tgbotapi.NewMessage(update.Message.Chat.ID, "No active games"))
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
			sendMessage(bot, msg)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
