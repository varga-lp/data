package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/varga-lp/data/books"
)

func main() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	bookChan := make(chan books.Book)
	closeChan := make(chan struct{})
	errChan := make(chan error)

	st := books.NewStreamer("BTCUSDT", bookChan, errChan, closeChan)
	if err := st.Dial(); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case book := <-bookChan:
				log.Println(book)
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
