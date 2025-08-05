// Package handlers is a nice package
package handlers

import (
	"fmt"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/odlev/tg-weather-bot/internal/openweather"
	"go.uber.org/zap"
)

type Handler struct {
	bot *tgbot.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbot.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot: bot,
		owClient: owClient,
	}
}

func (h *Handler) Start(log *zap.SugaredLogger) {
	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.HandleUpdate(update, log)
	}
}

func (h *Handler) HandleUpdate(update tgbot.Update, log *zap.SugaredLogger) {
	if update.Message != nil { // if we got a message
			log.Infof("[%s]: %s", update.Message.From.UserName, update.Message.Text)

			coordinates, err := h.owClient.GetCoordinates(update.Message.Text)
			if err != nil {
				log.Error("failed to get coordinates", zap.Error(err))
				h.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Не смог получить координаты"))

				return
			}
			weather, err := h.owClient.GetWeather(coordinates.Name, coordinates.Lat, coordinates.Lon)
			if err != nil {
				log.Error("failed to get weather", zap.Error(err))
				h.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Не смог получить погоду"))

				return
			}

			msg := tgbot.NewMessage(
				update.Message.Chat.ID, 
				fmt.Sprintf("%s\nПогода: %s\nПодробнее: %s\nТемпература: %.1f°C\nОщущается как: %.1f°C\nСкорость ветра: %.1f м/с", coordinates.Name, weather.Weather, weather.Description, weather.Temp, weather.FeelsLike, weather.Speed))
			msg.ReplyToMessageID = update.Message.MessageID

			h.bot.Send(msg)
		}
}

/* func (h *Handler) ButtonHandler(update tgbot.Update, log *zap.SugaredLogger) {
	replyKeyboard := tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("Хочу посмотреть погоду"),
			tgbot.NewKeyboardButton("Хочу запланировать отправку погоды"),
		),
	)

	switch update.Message.Text {
	case "Хочу посмотреть погоду":
		h.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Жду город"))
	case "Хочу запланировать отправку погоды":
	}
} 

func (h *Handler) WeatherMsg(update tgbot.Update, log *zap.SugaredLogger) {
	
} */
