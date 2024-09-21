package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var currentGameId = 0

type Game struct {
	Id       int
	Date     string
	Time     string
	Duration string
	Place    string
	Level    string
	Players  []string
}

func (g Game) String() string {
	return fmt.Sprint("Id: ", g.Id, ", date: ", g.Date, ", time: ", g.Time, ", duration: ", g.Duration, ", place: ", g.Place, ", level: ", g.Level, ", players: ", g.Players)
}

var activeGames []Game

type NewGameState int

const (
	NotStarted NewGameState = iota
	Started
	AwaitDate
	AwaitTime
	AwaitDuration
	AwaitPlace
	AwaitLevel
)

var gameState NewGameState = NotStarted

func main() {
	token := os.Getenv("PADEL_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if handleNewGameMessage(bot, update) {
			continue
		}
		if handleCommand(bot, update) {
			continue
		}
	}
}

func handleNewGameMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	text, ok := transitionGameState(update.Message.Text)
	if !ok {
		return false
	}
	sendMessage(bot, tgbotapi.NewMessage(update.Message.Chat.ID, text))
	return true
}

func transitionGameState(userInput string) (string, bool) {
	if gameState == NotStarted {
		return "", false
	}
	if gameState == Started {
		activeGames = append(activeGames, Game{Id: currentGameId})
		currentGameId++
	}
	game := &activeGames[len(activeGames)-1]
	switch gameState {
	case Started:
		gameState = AwaitDate
		return "Please enter the game date", true
	case AwaitDate:
		gameState = AwaitTime
		game.Date = userInput
		return "Please enter the game time", true
	case AwaitTime:
		gameState = AwaitDuration
		game.Time = userInput
		return "Please enter the game duration", true
	case AwaitDuration:
		gameState = AwaitPlace
		game.Duration = userInput
		return "Please enter the game place", true
	case AwaitPlace:
		gameState = AwaitLevel
		game.Place = userInput
		return "Please enter the game level", true
	case AwaitLevel:
		gameState = NotStarted
		game.Level = userInput
		return "Thank you, the game has been created", true
	}
	return "", false
}

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	if strings.HasPrefix(update.Message.Command(), "join") {
		handleJoinGame(bot, update)
		return true
	}
	switch update.Message.Command() {
	case "help":
		handleHelp(bot, update)
	case "new":
		handleNewGame(bot, update)
	case "games":
		handleActiveGames(bot, update)
	}
	return true
}

func handleNewGame(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	gameState = Started
	text, _ := transitionGameState("")
	msg.Text = text
	sendMessage(bot, msg)
}

func handleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Supported commands: /help to show this message, /new to create a new game, /games to show all active games"
	sendMessage(bot, msg)
}

func handleJoinGame(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if len(activeGames) == 0 {
		return
	}
	gameId, _ := strconv.Atoi(update.Message.Command()[4:])
	game := &activeGames[slices.IndexFunc(activeGames, func(g Game) bool { return g.Id == gameId })]
	game.Players = append(game.Players, "Test player")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Joined the game")
	sendMessage(bot, msg)
}

func handleActiveGames(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if len(activeGames) == 0 {
		msg.Text = "No active games"
	} else {
		var strs []string
		for _, g := range activeGames {
			gameStr := fmt.Sprint("Date: ", g.Date, "\nTime: ", g.Time, "\nDuration: ", g.Duration, "\nLevel: ", g.Level, "\nPlace: ", g.Place, "\nJoin the game: /join", g.Id)
			strs = append(strs, gameStr)
		}
		msg.Text = strings.Join(strs, "\n\n")
	}
	sendMessage(bot, msg)
}

func sendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
