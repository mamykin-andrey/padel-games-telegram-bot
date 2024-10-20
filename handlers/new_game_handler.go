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
	if !isReplyToTheBot(h.bot, update) || (update.Message != nil && (!isUserCreatingGame(update.Message.From.ID) || update.Message.IsCommand())) {
		return false
	}
	return transitionGameState(h.bot, update)
}

func isUserCreatingGame(userId int64) bool {
	_, exists := userGameStates[userId]
	return exists
}

func transitionGameState(bot shared.BotAPI, update tgbotapi.Update) bool {
	var input string
	var userId int64
	var userMessageId int
	var messageId int
	var chatId int64
	if update.Message != nil {
		input = update.Message.Text
		userId = update.Message.From.ID
		userMessageId = update.Message.MessageID
		chatId = update.Message.Chat.ID
		if update.Message.ReplyToMessage != nil {
			messageId = update.Message.ReplyToMessage.MessageID
		} else {
			messageId = -1
		}
	} else {
		input = update.CallbackQuery.Data
		userId = update.CallbackQuery.From.ID
		userMessageId = update.CallbackQuery.Message.MessageID
		chatId = update.CallbackQuery.Message.Chat.ID
		messageId = update.CallbackQuery.Message.MessageID
	}
	gameState := userGameStates[userId]
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
		bot.EditMessage(chatId, messageId, "Please enter the game time")
		bot.DeleteMessage(chatId, userMessageId)
		return true
	case AwaitTime:
		userGameStates[userId] = AwaitDuration
		game.Time = input
		bot.EditMessage(chatId, messageId, "Please enter the game duration")
		bot.DeleteMessage(chatId, userMessageId)
		return true
	case AwaitDuration:
		userGameStates[userId] = AwaitPlace
		game.Duration = input
		bot.EditMessage(chatId, messageId, "Please enter the game place")
		bot.DeleteMessage(chatId, userMessageId)
		return true
	case AwaitPlace:
		userGameStates[userId] = AwaitPlayers
		game.Place = input
		spotsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("1", "1"),
				tgbotapi.NewInlineKeyboardButtonData("2", "2"),
				tgbotapi.NewInlineKeyboardButtonData("3", "3"),
			),
		)
		bot.EditMessageTextAndMarkup(chatId, messageId, "Please enter how many spots you have", spotsKeyboard)
		bot.DeleteMessage(chatId, userMessageId)
		return true
	case AwaitPlayers:
		userGameStates[userId] = AwaitLevel
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > 3 {
			bot.EditMessage(chatId, messageId, "Please enter a correct number")
			bot.DeleteMessage(chatId, userMessageId)
			return true
		}
		game.NumberOfSpots = num
		levelKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("low begin", "low begin"),
				tgbotapi.NewInlineKeyboardButtonData("mid begin", "mid begin"),
				tgbotapi.NewInlineKeyboardButtonData("high begin", "high begin"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("low int", "low int"),
				tgbotapi.NewInlineKeyboardButtonData("mid int", "mid int"),
				tgbotapi.NewInlineKeyboardButtonData("high int", "high int"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("low adv", "low adv"),
				tgbotapi.NewInlineKeyboardButtonData("mid adv", "mid adv"),
				tgbotapi.NewInlineKeyboardButtonData("high adv", "high adv"),
			),
		)
		bot.EditMessageTextAndMarkup(chatId, messageId, "Please enter the game level", levelKeyboard)
		return true
	case AwaitLevel:
		userGameStates[userId] = NotStarted
		game.Level = input
		game.IsPublished = true
		game.Players = append(game.Players, fmt.Sprint("@", update.CallbackQuery.From.UserName))
		game.CreatorId = update.CallbackQuery.From.ID
		bot.DeleteMessage(chatId, messageId)
		NewActiveGamesCommandHandler(bot).ShowAllGames(chatId)
		return true
	}
	return false
}

func isReplyToTheBot(bot shared.BotAPI, update tgbotapi.Update) bool {
	// Special case for handling inline keyboard
	if update.CallbackQuery != nil {
		return true
	}
	if update.Message.ReplyToMessage == nil {
		return false
	}
	return bot.ID() == update.Message.ReplyToMessage.From.ID
}
