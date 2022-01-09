package currency

import (
	"encoding/json"
	"go.tgbot.crypto-currency-checker/internal/entities"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func GetCryptoCurrencyFromRemoteAPI(ApiUrl string, logger *zap.Logger) (entities.CryptoCurrencies, error) {
	resp, err := http.Get(ApiUrl)
	if err != nil {
		logger.Info("Не удалось сделать http запрос")
	}

	defer resp.Body.Close()

	reqBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Info("Не удалось прочесть Body ответа от API")
	}

	curr := entities.CryptoCurrencies{}

	err = json.Unmarshal(reqBody, &curr)
	if err != nil {
		logger.Info("Не удалось распарсить курс валют")
	}

	return curr, err
}
