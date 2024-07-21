package entity

import (
	"container/heap"
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

type BidsHeap struct {
	Data [][]PriceSize
}
type AsksHeap struct {
	Data [][]PriceSize
}

type teststruct struct {
	member1 int
	member2 int
}

func (h BidsHeap) Len() int { return len(h.Data) }
func (h AsksHeap) Len() int { return len(h.Data) }

func (h BidsHeap) Less(i, j int) bool {
	if len(h.Data[j]) == 0 {
		return true
	} else if len(h.Data[i]) == 0 {
		return false
	} else if h.Data[i][0].Price > h.Data[j][0].Price {
		return true
	} else {
		return false
	}
}
func (h AsksHeap) Less(i, j int) bool {
	if len(h.Data[j]) == 0 {
		return false
	} else if len(h.Data[i]) == 0 {
		return true
	} else if h.Data[i][0].Price < h.Data[j][0].Price {
		return true
	} else {
		return false
	}
}
func (h BidsHeap) Swap(i, j int) {
	h.Data[i], h.Data[j] = h.Data[j], h.Data[i]
}
func (h AsksHeap) Swap(i, j int) {
	h.Data[i], h.Data[j] = h.Data[j], h.Data[i]
}

func (h *BidsHeap) Push(x interface{}) {
	(*h).Data = append((*h).Data, x.([]PriceSize))
}
func (h *AsksHeap) Push(x interface{}) {
	(*h).Data = append((*h).Data, x.([]PriceSize))
}

func (h *BidsHeap) Pop() interface{} {
	old := (*h).Data
	n := len(old)
	x := old[n-1]
	(*h).Data = old[0 : n-1]
	return x
}
func (h *AsksHeap) Pop() interface{} {
	old := (*h).Data
	n := len(old)
	x := old[n-1]
	(*h).Data = old[0 : n-1]
	return x
}

func Aggregate(slice heap.Interface) []PriceSize {

	heap.Init(slice)
	aggregated := []PriceSize{}
	for slice.Len() != 0 {
		curr := heap.Pop(slice).([]PriceSize)

		if len(curr) == 0 {
			continue
		}
		heap.Push(slice, curr[1:])

		if len(aggregated) != 0 && aggregated[len(aggregated)-1].Price == curr[0].Price {
			aggregated[len(aggregated)-1].Size += curr[0].Size
		} else {
			aggregated = append(aggregated, curr[0])
		}
	}
	return aggregated
}
