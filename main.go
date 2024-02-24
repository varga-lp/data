package main

import (
	"log"
	"time"

	"github.com/varga-lp/data/klines"
)

func main() {
	log.Println(klines.Fetch("BTCUSDT", time.Now().UnixMilli()))
}
