package klines

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/varga-lp/data/config"
)

type klineResp struct {
	K struct {
		OpenTime       int64  `json:"t"`
		Open           string `json:"o"`
		High           string `json:"h"`
		Low            string `json:"l"`
		Close          string `json:"c"`
		Volume         string `json:"v"`
		TakerBuyVolume string `json:"V"`
		NumberOfTrades int64  `json:"n"`
		CloseTime      int64  `json:"T"`
		IsFinal        bool   `json:"x"`
		Omit           int64  `json:"L"` // this is to avoid case-insensitive unmarshalling
	} `json:"k"`
}

func (kr *klineResp) toKline() (Kline, error) {
	open, err := strconv.ParseFloat(kr.K.Open, 64)
	if err != nil {
		return Kline{}, err
	}
	high, err := strconv.ParseFloat(kr.K.High, 64)
	if err != nil {
		return Kline{}, err
	}
	low, err := strconv.ParseFloat(kr.K.Low, 64)
	if err != nil {
		return Kline{}, err
	}
	close, err := strconv.ParseFloat(kr.K.Close, 64)
	if err != nil {
		return Kline{}, err
	}
	volume, err := strconv.ParseFloat(kr.K.Volume, 64)
	if err != nil {
		return Kline{}, err
	}
	takerBuyVolume, err := strconv.ParseFloat(kr.K.TakerBuyVolume, 64)
	if err != nil {
		return Kline{}, err
	}

	return Kline{
		OpenTime:       kr.K.OpenTime,
		Open:           open,
		High:           high,
		Low:            low,
		Close:          close,
		Volume:         volume,
		TakerBuyVolume: takerBuyVolume,
		NumberOfTrades: kr.K.NumberOfTrades,
		CloseTime:      kr.K.CloseTime,
		IsFinal:        kr.K.IsFinal,
	}, nil
}

type Streamer struct {
	symbol     string
	streamName string
	updateChan chan<- Kline
	errorChan  chan<- error
	closeChan  <-chan struct{}
}

func NewStreamer(symbol string, updateChan chan<- Kline, errChan chan<- error, closeChan <-chan struct{}) *Streamer {
	symL := strings.ToLower(symbol)

	return &Streamer{
		symbol:     symbol,
		streamName: fmt.Sprintf("ws/%s@kline_%s", symL, KlineInterval),
		updateChan: updateChan,
		errorChan:  errChan,
		closeChan:  closeChan,
	}
}

func (s *Streamer) Dial() error {
	u := url.URL{Scheme: "wss", Host: config.MarketStreamEP(), Path: s.streamName}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-s.closeChan:
				c.Close()
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					s.errorChan <- err
					return
				}
				var kr klineResp
				if err := json.Unmarshal(message, &kr); err != nil {
					s.errorChan <- err
					return
				}
				kline, err := kr.toKline()
				if err != nil {
					s.errorChan <- err
					return
				}
				s.updateChan <- kline
			}
		}
	}()
	return nil
}
