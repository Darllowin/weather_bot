package handler

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/Darllowin/weather_bot/clients/openweather"
	"github.com/Darllowin/weather_bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userRepository interface {
	GetUserCity(ctx context.Context, userID int64) (string, error)
	CreateUser(ctx context.Context, userID int64) error
	UpdateCity(ctx context.Context, userID int64, city string) error
	GetUser(ctx context.Context, userID int64) (*models.User, error)
}

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
	userRepo userRepository
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient, userRepo userRepository) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
		userRepo: userRepo,
	}
}

func (h *Handler) handlerUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	ctx := context.Background()

	if update.Message.IsCommand() {
		err := h.ensureUser(ctx, update)
		if err != nil {
			log.Printf("error: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error")
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}

		switch update.Message.Command() {
		case "city":
			h.handleSetCity(ctx, update)
			return
		case "weather":
			h.handleSendWeather(ctx, update)
			return
		default:
			h.handleUnknownCommand(update)
			return
		}
	}
}

func (h *Handler) handleUnknownCommand(update tgbotapi.Update) {
	log.Printf("Unknown command [%s] %s", update.Message.From.UserName, update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Command not found")
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleSetCity(ctx context.Context, update tgbotapi.Update) {
	city := update.Message.CommandArguments()
	err := h.userRepo.UpdateCity(ctx, update.Message.From.ID, city)
	if err != nil {
		log.Printf("error: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("City: %s saved", city))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleSendWeather(ctx context.Context, update tgbotapi.Update) {
	city, err := h.userRepo.GetUserCity(ctx, update.Message.From.ID)
	if err != nil {
		log.Printf("error: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	if city == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Specify the city using the /city {city} command")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	coordinates, err := h.owClient.Coordinates(city)
	if err != nil {
		log.Printf("error: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "failed to get coordinates")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Printf("error: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "failed to get weather")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("%s: %d°C", city, int(math.Round(weather.Temp))))
	msg.ReplyToMessageID = update.Message.MessageID

	h.bot.Send(msg)
}

func (h *Handler) ensureUser(ctx context.Context, update tgbotapi.Update) error {
	user, err := h.userRepo.GetUser(ctx, update.Message.From.ID)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	if user == nil {
		err := h.userRepo.CreateUser(ctx, update.Message.From.ID)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
	}
	return nil
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handlerUpdate(update)
	}
}
