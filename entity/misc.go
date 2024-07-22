package entity

import (
	"time"
)

type Orderbook struct {
	Bids        []PriceSize
	Asks        []PriceSize
	Ts_original time.Time
	Ts_received time.Time
	Ts_updated  time.Time
}

type OrderBookMsg struct {
	Exchange  string
	Ticker    string
	Orderbook Orderbook
}
type PriceSize struct {
	Price float32
	Size  float32 `mapstructure:"qty"`
}

var TickerMap = map[string]string{"BTC/USDT": "btcusdt", "ETH/USDT": "ethusdt"}

var CentralizedOrderBooks = map[string]map[string]Orderbook{}

var ChannelMap = map[string]chan OrderBookMsg{"btcusdt": make(chan OrderBookMsg), "ethusdt": make(chan OrderBookMsg)}
