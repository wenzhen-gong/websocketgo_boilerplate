package util

import (
	"container/heap"
	"sort"
	"wz/entity"
)

func InsertBids(existing *[]entity.PriceSize, update []entity.PriceSize) {

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
					*existing = append((*existing)[:idx], append([]entity.PriceSize{elem}, (*existing)[idx+1:]...)...)
				} else {
					*existing = append((*existing)[:idx], (*existing)[idx+1:]...)
				}
			} else if (*existing)[idx].Price < elem.Price {
				*existing = append((*existing)[:idx], append([]entity.PriceSize{elem}, (*existing)[idx:]...)...)
			}
		}

	}
}

func InsertAsks(existing *[]entity.PriceSize, update []entity.PriceSize) {

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
					*existing = append((*existing)[:idx], append([]entity.PriceSize{elem}, (*existing)[idx+1:]...)...)
				} else {
					*existing = append((*existing)[:idx], (*existing)[idx+1:]...)
				}
			} else if (*existing)[idx].Price > elem.Price {
				*existing = append((*existing)[:idx], append([]entity.PriceSize{elem}, (*existing)[idx:]...)...)
			}
		}
	}
}

type BidsHeap struct {
	Data [][]entity.PriceSize
}
type AsksHeap struct {
	Data [][]entity.PriceSize
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
	(*h).Data = append((*h).Data, x.([]entity.PriceSize))
}
func (h *AsksHeap) Push(x interface{}) {
	(*h).Data = append((*h).Data, x.([]entity.PriceSize))
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

func Aggregate(slice heap.Interface) []entity.PriceSize {

	heap.Init(slice)
	aggregated := []entity.PriceSize{}
	for slice.Len() != 0 {
		curr := heap.Pop(slice).([]entity.PriceSize)

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
