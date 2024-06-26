package entity

type Orderbook struct {
	bids        []PrizeSize
	asks        []PrizeSize
	ts_original int
	ts_received int
	ts_updated  int
}

type PrizeSize struct {
	price float32
	size  float32
}
