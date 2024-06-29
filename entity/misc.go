package entity

import "time"

type Orderbook struct {
	Bids        []PriceSize
	Asks        []PriceSize
	Ts_original time.Time
	Ts_received time.Time
	Ts_updated  time.Time
}

type PriceSize struct {
	Price float32
	Size  float32 `mapstructure:"qty"`
}

var TickerMap = map[string]string{"BTC/USDT": "btcusdt", "ETH/USDT": "ethusdt"}

type OrderBookMsg struct {
	Exchange  string
	Ticker    string
	Orderbook Orderbook
}
