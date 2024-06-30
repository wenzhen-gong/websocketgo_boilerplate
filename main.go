package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
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

	// // get all possible tickers
	// v := []string{}
	// for _, value := range entity.TickerMap {
	// 	v = append(v, value)
	// }
	// slices.Sort(v)
	// slices.Compact(v)

	go Kraken(entity.ChannelMap)
	go aggregator_btcusdt(entity.ChannelMap["btcusdt"])
	go aggregator_ethusdt(entity.ChannelMap["ethusdt"])

	<-interrupt
	log.Println("User interruption")
	return

}

func Kraken(channelMap map[string]chan entity.OrderBookMsg) {
	sub := map[string]interface{}{"method": "subscribe", "params": map[string]interface{}{"channel": "book",
		"symbol": []string{"BTC/USDT", "ETH/USDT"},
		"depth":  1000}}

	e := exchange.New("kraken", "wss://ws.kraken.com/v2", sub)
	e.Connect()
	defer e.Connection.Close()
	e.SendSubMsg()
	e.ReceiveMsg(kraken.ParseKrakenData, channelMap)
}

func aggregator_btcusdt(ch <-chan entity.OrderBookMsg) {
	for {
		fmt.Println("aggregator btcusdt received: ", <-ch)
	}
}

func aggregator_ethusdt(ch <-chan entity.OrderBookMsg) {
	for {
		fmt.Println("aggregator ethusdt received: ", <-ch)
	}
}
