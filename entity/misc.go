package entity

import (
	"sort"
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

type PSheap [][]PriceSize

func (h PSheap) Len() int { return len(h) }

func (h PSheap) Less(i, j int) bool {
	if len(h[j]) == 0 {
		return true
	} else if len(h[i]) == 0 {
		return false
	} else if h[i][0].Price > h[j][0].Price {
		return true
	} else {
		return false
	}

}
func (h PSheap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *PSheap) Push(x interface{}) {
	*h = append(*h, x.([]PriceSize))
}

func (h *PSheap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

var TickerMap = map[string]string{"BTC/USDT": "btcusdt", "ETH/USDT": "ethusdt"}

var CentralizedOrderBooks = map[string]map[string]Orderbook{}

var ChannelMap = map[string]chan OrderBookMsg{"btcusdt": make(chan OrderBookMsg), "ethusdt": make(chan OrderBookMsg)}

func InsertBids(existing *[]PriceSize, update []PriceSize) {

	for _, elem := range update {
		idx := sort.Search(len(*existing), func(i int) bool {
			return (*existing)[i].Price <= elem.Price
		})
		if idx == len(*existing) {
			if elem.Size != 0 {
				*existing = append((*existing), elem)
			}
		} else {

			if (*existing)[idx].Price == elem.Price {
				if elem.Size != 0 {
					*existing = append((*existing)[:idx], append([]PriceSize{elem}, (*existing)[idx+1:]...)...)
				} else {
					*existing = append((*existing)[:idx], (*existing)[idx+1:]...)
				}
			} else if (*existing)[idx].Price < elem.Price {
				*existing = append((*existing)[:idx], append([]PriceSize{elem}, (*existing)[idx:]...)...)
			}
		}

	}
}

func InsertAsks(existing *[]PriceSize, update []PriceSize) {

	for _, elem := range update {
		idx := sort.Search(len(*existing), func(i int) bool {
			return (*existing)[i].Price >= elem.Price
		})
		if idx == len(*existing) {
			if elem.Size != 0 {
				*existing = append((*existing), elem)
			}
		} else {

			if (*existing)[idx].Price == elem.Price {
				if elem.Size != 0 {
					*existing = append((*existing)[:idx], append([]PriceSize{elem}, (*existing)[idx+1:]...)...)
				} else {
					*existing = append((*existing)[:idx], (*existing)[idx+1:]...)
				}
			} else if (*existing)[idx].Price > elem.Price {
				*existing = append((*existing)[:idx], append([]PriceSize{elem}, (*existing)[idx:]...)...)
			}
		}
	}
}
