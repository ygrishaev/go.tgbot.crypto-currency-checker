package main

import (
	"fmt"
	"go.tgbot.crypto-currency-checker/internal/config"
	"go.tgbot.crypto-currency-checker/internal/entities"
	"go.tgbot.crypto-currency-checker/internal/services/currency"
	"log"
	"math/rand"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func main() {
	// Делает рандом более рандомным
	rand.Seed(time.Now().UnixNano())

	appConfig, err := config.New("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to read app config: %v", err)
	}

	logger, _ := zap.NewProduction()
	logger.Info("Application start")

	curr := currency.GetCryptoCurrencyFromRemoteAPI(appConfig.CoincapApiUrl, logger)

	bot, err := tgbotapi.NewBotAPI(appConfig.TelegramAPIToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case entities.CryptoCurrencyBitcoin:
			price, _ := strconv.ParseFloat(curr.Data[0].PriceUsd, 64)
			msg.Text = fmt.Sprintf("Цена 1 BTC: %.4f usd 💰", price)
		case entities.CryptoCurrencyEthereum:
			price, _ := strconv.ParseFloat(curr.Data[1].PriceUsd, 64)
			msg.Text = fmt.Sprintf("Цена 1 ETH: %.4f usd 💰", price)
		default:
			msg.Text = "Введи /" +
				entities.CryptoCurrencyBitcoin + " или /" +
				entities.CryptoCurrencyEthereum + ", чтобы узнать текущую цена на криптовалюту"
		}

		if _, err := bot.Send(msg); err != nil {
			logger.Info("Не удалось послать ответ в telegram")
		}
	}

}
