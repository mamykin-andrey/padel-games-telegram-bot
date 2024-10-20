package handlers

import (
	"fmt"
	"main/models"
	"main/shared"
	"slices"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type JoinGameCommandHandler struct {
	bot shared.BotAPI
}

func NewJoinGameCommandHandler(bot shared.BotAPI) *JoinGameCommandHandler {
	return &JoinGameCommandHandler{bot: bot}
}

func (h *JoinGameCommandHandler) HandleCommand(update tgbotapi.Update) bool {
	command := update.CallbackQuery.Data
	fromUserName := update.CallbackQuery.From.UserName
	chatId := update.CallbackQuery.Message.Chat.ID
	if len(shared.Games) == 0 {
		return true
	}
	gameId, _ := strconv.Atoi(command[4:])
	game := &shared.Games[slices.IndexFunc(shared.Games, func(g models.Game) bool { return g.Id == gameId })]
	userName := fmt.Sprint("@", fromUserName)
	if isPlayerInGame(*game, userName) || game.IsFull() {
		h.bot.SendMessage(tgbotapi.NewMessage(chatId, "Already joined the game"))
		return true
	}
	game.Players = append(game.Players, userName)
	h.bot.SendMessage(tgbotapi.NewMessage(chatId, "Joined the game"))
	NewActiveGamesCommandHandler(h.bot).ShowAllGames(chatId)
	return true
}

func isPlayerInGame(game models.Game, userName string) bool {
	for _, p := range game.Players {
		if p == userName {
			return true
		}
	}
	return false
}
