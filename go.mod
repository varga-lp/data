module github.com/varga-lp/data

go 1.21.6

require github.com/gorilla/websocket v1.5.1

require golang.org/x/net v0.17.0 // indirect

replace github.com/varga-lp/data/symbols => ../symbols

replace github.com/varga-lp/data/klines => ../klines

replace github.com/varga-lp/data/config => ../config

replace github.com/varga-lp/data/books => ../books