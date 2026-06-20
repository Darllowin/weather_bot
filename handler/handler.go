package handler

import (
	"fmt"
	"log"
	"math"

	"github.com/Darllowin/weather_bot/clients/openweather"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
	}
}

func (h *Handler) handlerUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	err, coordinates := h.owClient.Coordinates(update.Message.Text)
	if err != nil {
		log.Printf("error: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "failed to get coordinates")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	err, weather := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Printf("error: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "failed to get weather")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Температура в %s: %d°C", update.Message.Text, int(math.Round(weather.Temp))))
	msg.ReplyToMessageID = update.Message.MessageID

	h.bot.Send(msg)
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handlerUpdate(update)
	}
}
