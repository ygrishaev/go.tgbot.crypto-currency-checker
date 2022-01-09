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
		case "Random":
			msg.Text = renderText(curr.Data[rand.Intn(len(curr.Data))])
		case entities.CryptoCurrencyBitcoin:
			msg.Text = renderText(curr.Data[0])
		case entities.CryptoCurrencyEthereum:
			msg.Text = renderText(curr.Data[1])
		case entities.CryptoCurrencyTether:
			msg.Text = renderText(curr.Data[2])
		case entities.CryptoCurrencySolana:
			msg.Text = renderText(curr.Data[4])
		case entities.CryptoCurrencyCardano:
			msg.Text = renderText(curr.Data[6])
		case entities.CryptoCurrencyDogecoin:
			msg.Text = renderText(curr.Data[11])
		default:
			msg.Text = fmt.Sprintf(
				"Список комманд, чтобы узнать текущую цену на криптовалюту:\n\n/%s\n/%s\n/%s\n/%s\n/%s\n/%s\n\n",
				entities.CryptoCurrencyBitcoin,
				entities.CryptoCurrencyEthereum,
				entities.CryptoCurrencyTether,
				entities.CryptoCurrencySolana,
				entities.CryptoCurrencyCardano,
				entities.CryptoCurrencyDogecoin,
			)

			msg.Text = msg.Text + fmt.Sprintf("Либо получите случайную монету через команду /Random")
		}

		if _, err := bot.Send(msg); err != nil {
			logger.Info("Не удалось послать ответ в telegram")
		}
	}

}

func renderText(curr *entities.CryptoCurrency) string {
	var result string
	price, _ := strconv.ParseFloat(curr.PriceUsd, 64)
	changePercent, _ := strconv.ParseFloat(curr.ChangePercent24Hr, 64)

	result = fmt.Sprintf("Цена %s по отношению к доллару США\n\n", curr.Name)
	result = result + fmt.Sprintf("1 %s = %.4f USD 💰\n", curr.Symbol, price)

	if changePercent >= 0 {
		result = result + "(+"
	} else {
		result = result + "("
	}
	result = result + fmt.Sprintf("%.2f%%󠀥 за 24 часа)\n\n", changePercent)
	result = result + fmt.Sprintf("Время обновления курса: %v\n\n", updatedAt.Format("2006-01-02 15:04"))

	if curr.Explorer != "" {
		result = result + fmt.Sprintf("Информация о монете: %s", curr.Explorer)
	}

	return result
}
