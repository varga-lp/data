package config

import "os"

var (
	env = os.Getenv("VARGA_ENV")
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

func MarketStreamEP() string {
	if EnvIsDev() {
		return "fstream.binancefuture.com"
	}
	return "fstream.binance.com"
}
