package handlers

import (
	"fmt"
	"main/models"
	"main/shared"
	"strconv"

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

type NewGameCommandHandler struct {
	bot shared.BotAPI
}

var currentGameId = 0
var userGameStates = make(map[int64]NewGameState)

func NewNewGameCommandHandler(bot shared.BotAPI) *NewGameCommandHandler {
	return &NewGameCommandHandler{bot: bot}
}

func (h *NewGameCommandHandler) HandleCommand(update tgbotapi.Update) bool {
	userGameStates[update.Message.From.ID] = Started
	transitionGameState(h.bot, update)
	return true
}

func (h *NewGameCommandHandler) HandleNewGameMessage(update tgbotapi.Update) bool {
	if !isReplyToTheBot(h.bot, update) || !isUserCreatingGame(update.Message.From.ID) || update.Message.IsCommand() {
		return false
	}
	return transitionGameState(h.bot, update)
}

func isUserCreatingGame(userId int64) bool {
	_, exists := userGameStates[userId]
	return exists
}

func transitionGameState(bot shared.BotAPI, update tgbotapi.Update) bool {
	userId := update.Message.From.ID
	input := update.Message.Text
	gameState := userGameStates[userId]
	userMessageId := update.Message.MessageID
	chatId := update.Message.Chat.ID
	if gameState == NotStarted {
		return false
	}
	if gameState == Started {
		shared.Games = append(shared.Games, models.Game{Id: currentGameId})
		currentGameId++
	}
	game := &shared.Games[len(shared.Games)-1]
	switch gameState {
	case Started:
		userGameStates[userId] = AwaitDate
		msg := tgbotapi.NewMessage(chatId, "Please enter the game date")
		bot.SendMessage(msg)
		return true
	case AwaitDate:
		userGameStates[userId] = AwaitTime
		game.Date = input
		bot.EditMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game time")
		bot.DeleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitTime:
		userGameStates[userId] = AwaitDuration
		game.Time = input
		bot.EditMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game duration")
		bot.DeleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitDuration:
		userGameStates[userId] = AwaitPlace
		game.Duration = input
		bot.EditMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game place")
		bot.DeleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitPlace:
		userGameStates[userId] = AwaitPlayers
		game.Place = input
		bot.EditMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter how many spots you have")
		bot.DeleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitPlayers:
		userGameStates[userId] = AwaitLevel
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > 3 {
			bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter a correct number"))
			return true
		}
		game.NumberOfSpots = num
		bot.EditMessage(chatId, update.Message.ReplyToMessage.MessageID, "Please enter the game level")
		bot.DeleteMessage(update.Message.Chat.ID, userMessageId)
		return true
	case AwaitLevel:
		userGameStates[userId] = NotStarted
		game.Level = input
		game.IsPublished = true
		game.Players = append(game.Players, fmt.Sprint("@", update.Message.From.UserName))
		game.CreatorId = update.Message.From.ID
		bot.DeleteMessage(update.Message.Chat.ID, userMessageId)
		bot.DeleteMessage(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
		return NewActiveGamesCommandHandler(bot).HandleCommand(update)
	}
	return false
}

func isReplyToTheBot(bot shared.BotAPI, update tgbotapi.Update) bool {
	if update.Message.ReplyToMessage == nil {
		return false
	}
	return bot.ID() == update.Message.ReplyToMessage.From.ID
}
