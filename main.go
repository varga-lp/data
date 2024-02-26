package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/varga-lp/data/klines"
)

func main() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	klineChan := make(chan klines.Kline)
	closeChan := make(chan struct{})
	errChan := make(chan error)

	st := klines.NewStreamer("BTCUSDT", klineChan, errChan, closeChan)
	if err := st.Dial(); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case kline := <-klineChan:
				log.Println(kline)
			case err := <-errChan:
				log.Fatal(err)
			}
		}
	}()

	<-shutdownCh
	go func() {
		closeChan <- struct{}{}
		os.Exit(0)
	}()
	<-time.After(5 * time.Second)
}
