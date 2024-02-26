package klines

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/varga-lp/data/config"
)

type Streamer struct {
	symbol     string
	streamName string
	updateChan chan Kline
	errorChan  chan error
	closeChan  chan struct{}
}

func NewStreamer(symbol string, updateChan chan Kline, errChan chan error, closeChan chan struct{}) *Streamer {
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
				mt, message, err := c.ReadMessage()
				if err != nil {
					s.errorChan <- err
					return
				}
				log.Println(mt, string(message))
			}
		}
	}()
	return nil
}
