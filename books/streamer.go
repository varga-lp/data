package books

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/varga-lp/data/config"
)

type Book struct {
	EventTime       int64  `json:"E"`
	TransactionTime int64  `json:"T"`
	BestBidPrice    string `json:"b"`
	BestBidQty      string `json:"B"`
	BestAskPrice    string `json:"a"`
	BestAskQty      string `json:"A"`
	Omit            string `json:"e"` // this is to avoid case-insensitive unmarshalling
}

type Streamer struct {
	symbol     string
	streamName string
	updateChan chan<- Book
	errChan    chan error
	closeChan  <-chan struct{}
}

func NewStreamer(symbol string, updateChan chan<- Book, closeChan <-chan struct{}) *Streamer {
	symL := strings.ToLower(symbol)

	return &Streamer{
		symbol:     symbol,
		streamName: fmt.Sprintf("ws/%s@bookTicker", symL),
		updateChan: updateChan,
		errChan:    make(chan error),
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
				return
			case err := <-s.errChan:
				defer c.Close()
				// reconnect on error
				defer s.Dial()

				log.Println(s.symbol, "book streamer error:", err)
				time.Sleep(config.WSReconnectBuffer())
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-s.closeChan:
				c.Close()
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					s.errChan <- err
					return
				}
				var bk Book
				if err := json.Unmarshal(message, &bk); err != nil {
					s.errChan <- err
					return
				}
				s.updateChan <- bk
			}
		}
	}()
	return nil
}
