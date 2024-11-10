package handlers

import (
	"fmt"
	"main/models"
	"main/shared"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ActiveGamesCommandHandler struct {
	bot shared.BotAPI
}

func NewActiveGamesCommandHandler(bot shared.BotAPI) *ActiveGamesCommandHandler {
	return &ActiveGamesCommandHandler{bot: bot}
}

func (h *ActiveGamesCommandHandler) HandleCommand(update tgbotapi.Update) {
	h.ShowAllGames(update.Message.Chat.ID)
}

func (h *ActiveGamesCommandHandler) ShowAllGames(chatId int64) {
	activeGames := make([]models.Game, 0)
	for _, g := range shared.State.Games() {
		gameDate, dateValid := shared.TryToParseDate(g.Date)
		if (!dateValid || !shared.IsDatePassed(gameDate)) && g.IsPublished {
			activeGames = append(activeGames, g)
		}
	}
	if len(activeGames) == 0 {
		h.bot.SendMessage(tgbotapi.NewMessage(chatId, "No active games"))
	}
	for _, g := range activeGames {
		gamePlayers := strings.Join(g.Players[:], ", ")
		gameStr := fmt.Sprint(
			"📅 Date: ", g.Date,
			"\n⏰ Time: ", g.Time,
			"\n⏲️ Duration: ", g.Duration,
			"\n📊 Level: ", g.Level,
			"\n📍 Location: ", g.Place,
			"\n🏋🏻‍♂️ Players: ", 4-g.NumberOfSpots, " + ", gamePlayers,
		)
		msg := tgbotapi.NewMessage(chatId, gameStr)
		actionsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Join", fmt.Sprint("join", g.Id)),
				tgbotapi.NewInlineKeyboardButtonData("Delete", fmt.Sprint("delete", g.Id)),
			),
		)
		msg.ReplyMarkup = actionsKeyboard
		h.bot.SendMessage(msg)
	}
}
