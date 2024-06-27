package exchange

import (
	"encoding/json"
	"fmt"
	"log"
	"wz/entity"
	"wz/utilities"

	"github.com/gorilla/websocket"
)

type Exchange struct {
	exchange   string
	address    string
	orderBooks *map[string]entity.Orderbook
	subMsg     map[string]interface{}
	Connection *websocket.Conn
}

func (exchange *Exchange) Connect() {

	// Create ws connection
	c, _, err := websocket.DefaultDialer.Dial(exchange.address, nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	// exchange.Connection = c
	defer c.Close()
	// send subscription message
	subMsg, _ := json.Marshal(exchange.subMsg)
	subMsgStr := []byte(subMsg)

	if connectionErr := c.WriteMessage(1, subMsgStr); connectionErr != nil {
		log.Println("Failed to send subcription message: ", connectionErr)
	}

	// listening for messages sent from ws server
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
		}

		var receivedMsg map[string]interface{}

		if err := json.Unmarshal([]byte(message), &receivedMsg); err != nil {
			panic(err)
		}
		// structure data and send to channel
		if utilities.ValidateKraken(receivedMsg) {
			fmt.Println("success")
		}

	}

}

func New(exchange string, address string, subMsg map[string]interface{}) *Exchange {
	return &Exchange{exchange, address, &map[string]entity.Orderbook{}, subMsg, nil}
}

func (exchange *Exchange) StuctureData() {

}
