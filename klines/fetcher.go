package klines

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/varga-lp/data/config"
)

const (
	KlineInterval = "3m"
)

type Kline struct {
	OpenTime       int64
	Open           float64
	High           float64
	Low            float64
	Close          float64
	Volume         float64
	TakerBuyVolume float64
	NumberOfTrades int64
	CloseTime      int64
	IsFinal        bool
}

func Fetch(symbol string, endTime int64) ([]Kline, error) {
	url := fmt.Sprintf("%s/fapi/v1/klines?symbol=%s&interval=%s&endTime=%d",
		config.FuturesEP(), symbol, KlineInterval, endTime)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var klines [][]interface{}
	err = json.Unmarshal(body, &klines)
	if err != nil {
		return nil, err
	}
	return parseKlines(klines)
}

func parseKlines(klines [][]interface{}) ([]Kline, error) {
	result := make([]Kline, 0, len(klines))

	for _, kl := range klines {
		if len(kl) != 12 {
			return nil, fmt.Errorf("invalid kline format: %v", kl)
		}

		floats := make([]float64, 0, 7)
		for _, i := range []int{1, 2, 3, 4, 5, 9} {
			f, ok := kl[i].(string)
			if !ok {
				return nil, fmt.Errorf("invalid kline field %d: %v", i, kl[i])
			}
			v, err := strconv.ParseFloat(f, 64)
			if err != nil {
				return nil, err
			}
			floats = append(floats, v)
		}

		openCloseNot := make([]int64, 0, 3)
		for _, i := range []int{0, 6, 8} {
			v, ok := kl[i].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid kline field %d: %v", i, kl[i])
			}
			openCloseNot = append(openCloseNot, int64(v))
		}

		result = append(result, Kline{
			OpenTime:       openCloseNot[0],
			Open:           floats[0],
			High:           floats[1],
			Low:            floats[2],
			Close:          floats[3],
			Volume:         floats[4],
			TakerBuyVolume: floats[5],
			NumberOfTrades: openCloseNot[2],
			CloseTime:      openCloseNot[1],
			IsFinal:        time.Now().UnixMilli() > openCloseNot[1],
		})
	}
	return result, nil
}
