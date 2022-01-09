package main

import (
	"fmt"
	"go.tgbot.crypto-currency-checker/internal/config"
	"go.tgbot.crypto-currency-checker/internal/entities"
	"go.tgbot.crypto-currency-checker/internal/services/currency"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var curr entities.CryptoCurrencies
var updatedAt time.Time = time.Now()

func main() {
	// Делает рандом более рандомным
	rand.Seed(time.Now().UnixNano())

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is a required variable")
	}

	appConfig, err := config.New(configPath)
	if err != nil {
		log.Fatalf("failed to read app config: %v", err)
	}

	logger, _ := zap.NewProduction()
	logger.Info("Application start")

	logger.Info("Отправка первичного запроса на получение курса криптовалют")
	BufferCurr, BufferTime, err := currency.GetCryptoCurrencyFromRemoteAPI(appConfig.CoincapApiUrl, logger)
	if err != nil {
		logger.Info("Ошибка запроса к API через Goroutine")
	} else {
		curr = BufferCurr
		updatedAt = BufferTime
	}

	go func() {
		for range time.Tick(time.Minute) {
			logger.Info("Отправка вторичного запроса на получение курса криптовалют")
			BufferCurr, BufferTime, err := currency.GetCryptoCurrencyFromRemoteAPI(appConfig.CoincapApiUrl, logger)
			if err != nil {
				logger.Info("Ошибка запроса к API через Goroutine")
			} else {
				curr = BufferCurr
				updatedAt = BufferTime
			}
		}
	}()

	bot, err := tgbotapi.NewBotAPI(appConfig.TelegramAPIToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case entities.CryptoCurrencyBitcoin:
			price, _ := strconv.ParseFloat(curr.Data[0].PriceUsd, 64)
			msg.Text = fmt.Sprintf("Цена 1 BTC: %.4f usd 💰", price)
			msg.Text = fmt.Sprintf(msg.Text+"\n\nВремя обновления курса: %s", updatedAt.Format("2006-01-02 15:04:05"))
		case entities.CryptoCurrencyEthereum:
			price, _ := strconv.ParseFloat(curr.Data[1].PriceUsd, 64)
			msg.Text = fmt.Sprintf("Цена 1 ETH: %.4f usd 💰", price)
			msg.Text = fmt.Sprintf(msg.Text+"\n\nВремя обновления курса: %s", updatedAt.Format("2006-01-02 15:04:05"))
		default:
			msg.Text = fmt.Sprintf(
				"Введите /%s или /%s, чтобы узнать текущую цену на криптовалюту",
				entities.CryptoCurrencyBitcoin,
				entities.CryptoCurrencyEthereum,
			)
		}

		if _, err := bot.Send(msg); err != nil {
			logger.Info("Не удалось послать ответ в telegram")
		}
	}

}
