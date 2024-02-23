package config

import "os"

const (
	HistoricKlineInterval      = "3m"
	HistoricKlineMinKlineCount = 1
	HistoricKlineMaxKlineCount = 31 * 24 * 60 / 3 // 31 days

	LiveDataKlineInterval = "3m"
	LiveDataMaxKlineCount = 250
	LiveDataMinKlineCount = 1
)

var (
	env         = os.Getenv("VDATA_ENV")
	DatabaseUrl = os.Getenv("VDATA_DATABASE_URL")
)

func EnvIsDev() bool {
	return !EnvIsProd()
}

func EnvIsProd() bool {
	return env == "prod"
}

func FuturesEP() string {
	if EnvIsDev() {
		return "https://testnet.binancefuture.com"
	}
	return "https://fapi.binance.com"
}
