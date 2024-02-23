package symbols

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/varga-lp/data/config"
)

const (
	targetStatus       = "TRADING"
	targetContractType = "PERPETUAL"
	targetQuoteAsset   = "USDT"
	targetFilter       = "PRICE_FILTER"
	targetAge          = 60
)

type exchangeInfoResp struct {
	Symbols []struct {
		Symbol            string   `json:"symbol"`
		Status            string   `json:"status"`
		ContractType      string   `json:"contractType"`
		QuoteAsset        string   `json:"quoteAsset"`
		UnderlyingSubType []string `json:"underlyingSubType"`
		OnboardDate       int64    `json:"onboardDate"`
		Filters           []struct {
			FilterType string `json:"filterType"`
			TickSize   string `json:"tickSize"`
		} `json:"filters"`
	} `json:"symbols"`
}

type symbol struct {
	Symbol   string
	TickSize float64
}

func fetch() ([]symbol, error) {
	url := config.FuturesEP() + "/fapi/v1/exchangeInfo"

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var respBody exchangeInfoResp
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return nil, err
	}
	return filterSymbols(respBody)
}

func filterSymbols(resp exchangeInfoResp) ([]symbol, error) {
	var symbols []symbol
	for _, sym := range resp.Symbols {
		if sym.Status == targetStatus && sym.ContractType == targetContractType &&
			sym.QuoteAsset == targetQuoteAsset && sym.OnboardDate < time.Now().AddDate(0, 0, -targetAge).UnixMilli() {
			for _, filter := range sym.Filters {
				if filter.FilterType == targetFilter {
					ts, err := strconv.ParseFloat(filter.TickSize, 64)
					if err != nil {
						return nil, err
					}

					symbols = append(symbols, symbol{
						Symbol:   sym.Symbol,
						TickSize: ts,
					})
					break
				}
			}
		}
	}
	return symbols, nil
}
