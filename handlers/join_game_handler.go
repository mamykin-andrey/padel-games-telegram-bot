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
	if len(shared.Games) == 0 {
		return true
	}
	gameId, _ := strconv.Atoi(update.Message.Command()[4:])
	game := &shared.Games[slices.IndexFunc(shared.Games, func(g models.Game) bool { return g.Id == gameId })]
	userName := fmt.Sprint("@", update.Message.From.UserName)
	if isPlayerInGame(*game, userName) || game.IsFull() {
		h.bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Already joined the game"))
		return true
	}
	game.Players = append(game.Players, userName)
	h.bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Joined the game"))
	return NewActiveGamesCommandHandler(h.bot).HandleCommand(update)
}

func isPlayerInGame(game models.Game, userName string) bool {
	for _, p := range game.Players {
		if p == userName {
			return true
		}
	}
	return false
}
