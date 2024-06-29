package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type commandStruct struct {
	Type        string   `json:"type"`
	Product_IDs []string `json:"product_ids"`
	Channels    []string `json:"channels"`
}

var addr = flag.String("addr", "ws-feed.exchange.coinbase.com", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// done := make(chan struct{})

	go func() {
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s, type: %s", message, websocket.FormatMessageType(mt))
		}
	}()

	command := &commandStruct{
		Type:        "subscribe",
		Product_IDs: []string{"BTC-USD"},
		Channels:    []string{"level2_batch"},
	}

	json, _ := json.Marshal(command)
	str := string(json)

	connectionErr := c.WriteMessage(1, []byte(str))
	if connectionErr != nil {
		log.Println("write:", connectionErr)
	}

	// select {
	// case <-interrupt:
	// 	log.Println("interrupttttttttttt")
	// 	return
	// }
	<-interrupt
	log.Println("interrupttttttttttt")
	return
}
