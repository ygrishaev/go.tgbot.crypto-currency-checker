package entities

type CryptoCurrency struct {
	ID                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
	Explorer          string `json:"explorer"`
}

type CryptoCurrencies struct {
	Data []*CryptoCurrency `json:"data"`
}

const (
	CryptoCurrencyBitcoin  string = "Bitcoin"
	CryptoCurrencyEthereum string = "Ethereum"
	CryptoCurrencyTether   string = "Tether"
	CryptoCurrencySolana   string = "Solana"
	CryptoCurrencyCardano  string = "Cardano"
	CryptoCurrencyDogecoin string = "Dogecoin"
)
