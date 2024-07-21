package main

import (
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

}

func aggregators(channelMap map[string]chan entity.OrderBookMsg) {
	for ticker, channel := range channelMap {
		go func() {
			for {
				// update CentralizedOrderBooks
				orderBookMsg := <-channel
				entity.CentralizedOrderBooks[ticker] = map[string]entity.Orderbook{orderBookMsg.Exchange: orderBookMsg.Orderbook}

				// construct bidsSlice and asksSlice to be stored and used as heap data structure
				bidsSlice := &entity.BidsHeap{Data: [][]entity.PriceSize{}}
				asksSlice := &entity.AsksHeap{Data: [][]entity.PriceSize{}}

				for _, orderBook := range entity.CentralizedOrderBooks[ticker] {
					(*bidsSlice).Data = append((*bidsSlice).Data, orderBook.Bids)
					(*asksSlice).Data = append((*asksSlice).Data, orderBook.Asks)

					// <--- SHOULD DELETE
					// Currently there is only Kraken in CentralizedOrderBooks, so Bids and Asks are added twice to make sure entity.Aggregate works properly
					(*bidsSlice).Data = append((*bidsSlice).Data, orderBook.Bids)
					(*asksSlice).Data = append((*asksSlice).Data, orderBook.Asks)
					// --->

				}

				// aggregate bids and asks (Feel free to uncomment and print aggregated bids and asks)
				// aggregatedBids := entity.Aggregate(bidsSlice)
				// aggregatedAsks := entity.Aggregate(asksSlice)

				// fmt.Println("aggregatedBids: ", aggregatedBids)
				// fmt.Println("aggregatedAsks: ", aggregatedAsks)

				entity.Aggregate(bidsSlice)
				entity.Aggregate(asksSlice)

				fmt.Println("Took ", time.Since(orderBookMsg.Orderbook.Ts_original), "to aggregate bids and asks (using Ts_original for worst scenario)")
			}
		}()
	}
}
