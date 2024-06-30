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

type PriceSize struct {
	Price float32
	Size  float32 `mapstructure:"qty"`
}

var TickerMap = map[string]string{"BTC/USDT": "btcusdt", "ETH/USDT": "ethusdt"}

var ChannelMap = map[string]chan OrderBookMsg{"btcusdt": make(chan OrderBookMsg), "ethusdt": make(chan OrderBookMsg)}

type OrderBookMsg struct {
	Exchange  string
	Ticker    string
	Orderbook Orderbook
}

// func UpdateOrderBooks(messageType string, ticker string, orderBook Orderbook, orderBooks map[string]*Orderbook, maxDepth int) {
// 	if messageType == "snapshot" {
// 		orderBook.Ts_updated = time.Now().UTC()
// 		orderBooks[ticker] = &orderBook
// 	} else if messageType == "update" {
// 		// insert update Bids and Asks into orderBooks, trim to maxDepth
// 		if orderBooks[ticker] != nil {

// 			insertBids(&orderBooks[ticker].Bids, orderBook.Bids)
// 			orderBooks[ticker].Bids = orderBooks[ticker].Bids[:int(math.Min(float64(len(orderBooks[ticker].Bids)), float64(maxDepth)))]
// 			insertAsks(&orderBooks[ticker].Asks, orderBook.Asks)
// 			orderBooks[ticker].Asks = orderBooks[ticker].Asks[:int(math.Min(float64(len(orderBooks[ticker].Asks)), float64(maxDepth)))]

// 			// update orderBooks's Ts_original, Ts_received, Ts_updated from update message and current time respectively
// 			orderBooks[ticker].Ts_original = orderBook.Ts_original
// 			orderBooks[ticker].Ts_received = orderBook.Ts_received
// 			orderBooks[ticker].Ts_updated = time.Now().UTC()
// 		}
// 	}
// }

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
