package main

import (
	"fmt"
	"slices"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type JoinGameCommandHandler struct {
	bot *tgbotapi.BotAPI
}

func NewJoinGameCommandHandler(bot *tgbotapi.BotAPI) *JoinGameCommandHandler {
	return &JoinGameCommandHandler{bot: bot}
}

func (h *JoinGameCommandHandler) HandleCommand(update tgbotapi.Update) bool {
	if len(games) == 0 {
		return true
	}
	gameId, _ := strconv.Atoi(update.Message.Command()[4:])
	game := &games[slices.IndexFunc(games, func(g Game) bool { return g.Id == gameId })]
	userName := fmt.Sprint("@", update.Message.From.UserName)
	if isPlayerInGame(*game, userName) || game.isFull() {
		sendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Already joined the game"))
		return true
	}
	game.Players = append(game.Players, userName)
	sendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Joined the game"))
	return true
}

func isPlayerInGame(game Game, userName string) bool {
	for _, p := range game.Players {
		if p == userName {
			return true
		}
	}
	return false
}
