package main

import (
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/odlev/tg-weather-bot/internal/handlers"
	"github.com/odlev/tg-weather-bot/internal/openweather"
	"github.com/odlev/tg-weather-bot/pkg/logger"
)
func main() {
	log := logger.SetupLogger()
	defer log.Sync()

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	

	bot, err := tgbot.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("error creating bot")
	}
	bot.Debug = false
	log.Infof("Authorized on account %s", bot.Self.UserName)
	
	owClient := openweather.New(os.Getenv("WEATHER_TOKEN"))

	botHandler := handlers.New(bot, owClient)
	
	botHandler.Start(log)
}

