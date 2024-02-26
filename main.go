package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/varga-lp/data/klines"
)

func main() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	closeChan := make(chan struct{})
	errChan := make(chan error)

	st := klines.NewStreamer("BTCUSDT", nil, errChan, closeChan)
	if err := st.Dial(); err != nil {
		panic(err)
	}

	<-shutdownCh
	go func() {
		closeChan <- struct{}{}
		os.Exit(0)
	}()
	<-time.After(5 * time.Second)
}
