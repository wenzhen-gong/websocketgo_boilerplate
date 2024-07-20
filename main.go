package main

import (
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"wz/entity"
	"wz/exchange"
	"wz/kraken"
)

func main() {
	// Overhead (setting up options for executable, register interrupt to receive notification of os.Interrupt signal)
	flag.Parse()
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go Kraken(entity.ChannelMap)
	aggregators(entity.ChannelMap)

	<-interrupt
	log.Println("User interruption")
	return

}

func Kraken(channelMap map[string]chan entity.OrderBookMsg) {
	for _, value := range entity.ChannelMap {
		defer close(value)
	}
	sub := map[string]interface{}{"method": "subscribe", "params": map[string]interface{}{"channel": "book",
		"symbol": []string{"BTC/USDT", "ETH/USDT"},
		"depth":  1000}}

	e := exchange.New("kraken", "wss://ws.kraken.com/v2", sub)
	e.Connect()
	defer e.Connection.Close()
	e.SendSubMsg()
	e.ReceiveMsg(kraken.ParseKrakenData, channelMap)
	return
}

func aggregators(channelMap map[string]chan entity.OrderBookMsg) {
	for ticker, channel := range channelMap {
		go func() {
			for {
				// update CentralizedOrderBooks
				orderBookMsg := <-channel
				entity.CentralizedOrderBooks[ticker] = map[string]entity.Orderbook{orderBookMsg.Exchange: orderBookMsg.Orderbook}

				// construct slice to be stored and used as heap data structure
				slice := &entity.PSheap{}
				for _, orderBook := range entity.CentralizedOrderBooks[ticker] {
					*slice = append(*slice, orderBook.Bids)
				}

				// aggregate
				heap.Init(slice)
				aggregatedBids := []entity.PriceSize{}
				for slice.Len() != 0 {
					curr := heap.Pop(slice).([]entity.PriceSize)
					if len(curr) == 0 {
						continue
					}
					heap.Push(slice, curr[1:])

					if len(aggregatedBids) != 0 && aggregatedBids[len(aggregatedBids)-1].Price == curr[0].Price {
						aggregatedBids[len(aggregatedBids)-1].Size += curr[0].Size
					} else {
						aggregatedBids = append(aggregatedBids, curr[0])
					}
				}
				fmt.Println("Took ", time.Since(orderBookMsg.Orderbook.Ts_original), "to aggregate (using Ts_original for worst scenario consideration)")
			}
		}()
	}
}
