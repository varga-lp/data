package config

import "os"

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
