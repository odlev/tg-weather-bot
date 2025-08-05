// Package handlers is a nice package
package handlers

import (
	"fmt"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/odlev/tg-weather-bot/internal/openweather"
	"go.uber.org/zap"
)

type UserState int

const (
	StateDefault UserState = iota
	StateWaitingForCity
	StateWaitingForInterval
)

var userStates = make(map[int64]UserState)

type Handler struct {
	bot      *tgbot.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbot.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:      bot,
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
	chatID := update.Message.Chat.ID

	if update.Message != nil { // if we got a message
		log.Infof("[%s]: %s", update.Message.From.UserName, update.Message.Text)

		switch userStates[chatID] {
		case StateWaitingForCity:
			h.sendWeatherByCity(update, log)

			userStates[chatID] = StateDefault
			h.sendMainMenu(chatID)

			return
		case StateWaitingForInterval:
			// TODO: код тикера
			userStates[chatID] = StateDefault
			h.sendMainMenu(chatID)

			return
		}

		switch update.Message.Text {
		case "/start":
			h.sendMainMenu(chatID)
		case "Хочу посмотреть погоду":
			userStates[chatID] = StateWaitingForCity
			msg := tgbot.NewMessage(chatID, "Введите название города:")
			msg.ReplyMarkup = tgbot.ReplyKeyboardRemove{RemoveKeyboard: true}
			h.bot.Send(msg)
		case "Хочу запланировать отправку погоды":
			userStates[chatID] = StateWaitingForInterval
			h.bot.Send(tgbot.NewMessage(chatID, "Введите интервал, через который вы будете получать оповещение о погоде:\n[не доделано]"))
			h.sendTimeMenu(chatID)
		case "Отмена":
			userStates[chatID] = StateDefault
			h.sendMainMenu(chatID)
		}

	}
}

func (h *Handler) sendMainMenu(chatID int64) {
	replyKeyboard := tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("Хочу посмотреть погоду"),
			tgbot.NewKeyboardButton("Хочу запланировать отправку погоды"),
		),
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("Отмена")),
	)

	msg := tgbot.NewMessage(chatID, "Главное меню:")
	msg.ReplyMarkup = replyKeyboard
	h.bot.Send(msg)
}

func (h *Handler) sendTimeMenu(chatID int64) {
	replyKeyboard := tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("10 секунд"),
			tgbot.NewKeyboardButton("1 минута"),
		), tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("1 час"),
			tgbot.NewKeyboardButton("24 часа"),
		),
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("Ввести свое время в минутах"),
			tgbot.NewKeyboardButton("Ввести свое время в часах"),
		),
	)
	msg := tgbot.NewMessage(chatID, "Установите интервал:")
	msg.ReplyMarkup = replyKeyboard
	h.bot.Send(msg)
}

func (h *Handler) sendWeatherByCity(update tgbot.Update, log *zap.SugaredLogger) {
	
	
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

/* func (h *Handler) startWeatherTicker(city string, log *zap.SugaredLogger, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	
	done := make(chan bool)

	go func() {
		time.Sleep(time.Second)
		close(done)
	}()

	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				return
			}

		case <-done:
			return
		}
	}
} */

/* func (h *Handler) getWeather(update tgbot.Update) tgbot.MessageConfig {
	coordinates, err := h.owClient.GetCoordinates(update.Message.Text)
	if err != nil {
		log.Error("failed to get coordinates", zap.Error(err))
		msg := tgbot.NewMessage(update.Message.Chat.ID, "Не смог получить координаты"))

		return msg
	}
	weather, err := h.owClient.GetWeather(coordinates.Name, coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Error("failed to get weather", zap.Error(err))
		msg := tgbot.NewMessage(update.Message.Chat.ID, "Не смог получить погоду"))

		return msg
	}
	return 
} */
