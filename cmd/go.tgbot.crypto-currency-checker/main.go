package main

import (
	"go.tgbot.crypto-currency-checker/internal/config"
	"go.tgbot.crypto-currency-checker/internal/entities"
	"go.tgbot.crypto-currency-checker/internal/services/currency"
	"go.tgbot.crypto-currency-checker/internal/services/response"
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var curr entities.CryptoCurrencies

func main() {

	// Делает рандом более рандомным
	rand.Seed(time.Now().UnixNano())

	// Прописываем при запуске как CONFIG_PATH=configs/config.yaml
	// Либо настраиваем в IDE переменную окружения
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is a required variable")
	}

	// appConfig - структура со всеми конфигами приложения
	appConfig, err := config.New(configPath)
	if err != nil {
		log.Fatalf("failed to read app config: %v", err)
	}

	logger, _ := zap.NewProduction()
	logger.Info("Application start")

	logger.Info("Отправка первичного запроса на получение курса криптовалют")
	curr, err = currency.GetCryptoCurrencyFromRemoteAPI(appConfig.CoincapApiUrl, logger)
	if err != nil {
		logger.Info("Ошибка запроса к API через Goroutine")
	}

	go func() {
		for range time.Tick(time.Minute) {
			logger.Info("Отправка вторичного запроса на получение курса криптовалют")
			curr, err = currency.GetCryptoCurrencyFromRemoteAPI(appConfig.CoincapApiUrl, logger)
			if err != nil {
				logger.Info("Ошибка запроса к API через Goroutine")
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
		case "Random":
			msg.Text = response.Render(curr.Data[rand.Intn(len(curr.Data))], curr.LastUpdate)
		case entities.CryptoCurrencyBitcoin:
			msg.Text = response.Render(curr.Data[0], curr.LastUpdate)
		case entities.CryptoCurrencyEthereum:
			msg.Text = response.Render(curr.Data[1], curr.LastUpdate)
		case entities.CryptoCurrencyTether:
			msg.Text = response.Render(curr.Data[2], curr.LastUpdate)
		case entities.CryptoCurrencySolana:
			msg.Text = response.Render(curr.Data[4], curr.LastUpdate)
		case entities.CryptoCurrencyCardano:
			msg.Text = response.Render(curr.Data[6], curr.LastUpdate)
		case entities.CryptoCurrencyDogecoin:
			msg.Text = response.Render(curr.Data[11], curr.LastUpdate)
		default:
			msg.Text = response.StartMessage()
		}

		if _, err := bot.Send(msg); err != nil {
			logger.Info("Не удалось послать ответ в telegram")
		}
	}

}
