package exchange

import (
	"encoding/json"
	"log"
	"wz/entity"

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
	json, _ := json.Marshal(exchange.subMsg)
	str := []byte(json)

	if connectionErr := c.WriteMessage(1, str); connectionErr != nil {
		log.Println("Failed to send subcription message: ", connectionErr)
	}

	// listening for messages sent from ws server
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}
		// structure data and send to channel

		log.Printf("recv: %s, type: %s", message, websocket.FormatMessageType(mt))
	}

}

func New(exchange string, address string, subMsg map[string]interface{}) *Exchange {
	return &Exchange{exchange, address, &map[string]entity.Orderbook{}, subMsg, nil}
}
