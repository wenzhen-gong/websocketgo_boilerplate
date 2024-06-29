package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"wz/exchange"
	"wz/kraken"
)

func main() {
	// Overhead (setting up options for executable, register interrupt to receive notification of os.Interrupt signal)
	flag.Parse()
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go Kraken()

	<-interrupt
	log.Println("User interruption")
	return

	// for {
	// 	select {
	// 	// case <-done:
	// 	// 	return
	// 	case <-interrupt:
	// 		log.Println("interrupt")

	// 		// Cleanly close the connection by sending a close message and then
	// 		// waiting (with timeout) for the server to close the connection.
	// 		err := e.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 		if err != nil {
	// 			log.Println("write close:", err)
	// 			return
	// 		}
	// 		// select {
	// 		// case <-done:
	// 		// case <-time.After(time.Second):
	// 		// }
	// 		return
	// 	}
	// }
}

func Kraken() {
	sub := map[string]interface{}{"method": "subscribe", "params": map[string]interface{}{"channel": "book",
		"symbol": []string{"BTC/USDT", "ETH/USDT"},
		"depth":  1000}}

	e := exchange.New("kraken", "wss://ws.kraken.com/v2", sub)
	e.Connect()
	defer e.Connection.Close()
	e.SendSubMsg()
	e.ReceiveMsg(kraken.ParseKrakenData)
}
