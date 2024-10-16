package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ActiveGamesCommandHandler struct {
	bot *tgbotapi.BotAPI
}

func NewActiveGamesCommandHandler(bot *tgbotapi.BotAPI) *ActiveGamesCommandHandler {
	return &ActiveGamesCommandHandler{bot: bot}
}

func (h *ActiveGamesCommandHandler) HandleCommand(update tgbotapi.Update) bool {
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
	return true
}
