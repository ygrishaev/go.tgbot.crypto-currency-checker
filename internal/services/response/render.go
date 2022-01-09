package response

import (
	"fmt"
	"go.tgbot.crypto-currency-checker/internal/entities"
	"strconv"
	"time"
)

func Render(curr *entities.CryptoCurrency, time time.Time) string {
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
	result = result + fmt.Sprintf("Время обновления курса: %v\n\n", time.Format("2006-01-02 15:04"))

	if curr.Explorer != "" {
		result = result + fmt.Sprintf("Информация о монете: %s", curr.Explorer)
	}

	return result
}

func StartMessage() string {
	message := fmt.Sprintf(
		"Список комманд, чтобы узнать текущую цену на криптовалюту:\n\n/%s\n/%s\n/%s\n/%s\n/%s\n/%s\n\n",
		entities.CryptoCurrencyBitcoin,
		entities.CryptoCurrencyEthereum,
		entities.CryptoCurrencyTether,
		entities.CryptoCurrencySolana,
		entities.CryptoCurrencyCardano,
		entities.CryptoCurrencyDogecoin,
	)

	return message + fmt.Sprintf("Либо получите случайную монету через команду /Random")
}
