package currency

import (
	"encoding/json"
	"go.tgbot.crypto-currency-checker/internal/entities"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

func GetCryptoCurrencyFromRemoteAPI(ApiUrl string, logger *zap.Logger) entities.CryptoCurrencies {
	resp, err := http.Get(ApiUrl)
	if err != nil {
		logger.Info("Не удалось сделать http запрос")
		log.Panic(err)
	}

	defer resp.Body.Close()

	curr := entities.CryptoCurrencies{}

	reqBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Info("Не удалось прочесть Body ответа от API")
		log.Panic(err)
	}

	err = json.Unmarshal(reqBody, &curr)
	if err != nil {
		logger.Info("Не удалось распарсить курс валют")
		log.Panic(err)
	}

	return curr
}
