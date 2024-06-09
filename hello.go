package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

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

	done := make(chan struct{})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		defer close(done)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s, type: %s", message, websocket.FormatMessageType(mt))
		}
	}()

	type commandStruct struct {
		Type        string   `json:"type"`
		Product_IDs []string `json:"product_ids"`
		Channels    []string `json:"channels"`
	}
	command := &commandStruct{
		Type:        "subscribe",
		Product_IDs: []string{"BTC-USD"},
		Channels:    []string{"level2_batch"},
	}

	json, _ := json.Marshal(command)
	str := string(json)

	fmt.Println(str)
	connectionErr := c.WriteMessage(1, []byte(str))
	if connectionErr != nil {
		log.Println("write:", connectionErr)
	}

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
