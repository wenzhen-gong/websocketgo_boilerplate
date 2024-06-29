package exchange

import (
	"encoding/json"
	"log"
	"math"
	"time"
	"wz/entity"

	"github.com/gorilla/websocket"
)

type Exchange struct {
	exchange   string
	address    string
	orderBooks map[string]*entity.Orderbook
	maxDepth   int
	subMsg     map[string]interface{}
	Connection *websocket.Conn
}

func (exchange *Exchange) Connect() {

	// Create ws connection
	c, _, err := websocket.DefaultDialer.Dial(exchange.address, nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	// defer c.Close()
	exchange.Connection = c
	// send subscription message
	// subMsg, _ := json.Marshal(exchange.subMsg)
	// subMsgStr := []byte(subMsg)

	// if connectionErr := c.WriteMessage(1, subMsgStr); connectionErr != nil {
	// 	log.Println("Failed to send subcription message: ", connectionErr)
	// }

	// listening for messages sent from ws server
	// for {
	// 	_, message, err := c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("Failed to read message:", err)
	// 	}

	// 	var receivedMsg map[string]interface{}

	// 	if err := json.Unmarshal([]byte(message), &receivedMsg); err != nil {
	// 		panic(err)
	// 	}
	// 	// structure data and send to channel
	// 	if utilities.ValidateKraken(receivedMsg) {
	// 		fmt.Println("success")
	// 	}

	// }

}

func New(exchange string, address string, subMsg map[string]interface{}) *Exchange {
	return &Exchange{exchange, address, map[string]*entity.Orderbook{}, 2000, subMsg, nil}
}

func (exchange *Exchange) SendSubMsg() {
	// fmt.Println(exchange)

	subMsg, _ := json.Marshal(exchange.subMsg)
	subMsgStr := []byte(subMsg)

	if connectionErr := exchange.Connection.WriteMessage(1, subMsgStr); connectionErr != nil {
		log.Println("Failed to send subcription message: ", connectionErr)
	}
}

func (exchange *Exchange) ReceiveMsg(parseData func(map[string]interface{}, int) (string, string, entity.Orderbook)) {
	for {
		_, message, err := exchange.Connection.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
		}
		// Parse received JSON from exchange ws server into Go map
		var receivedMsg map[string]interface{}
		if err := json.Unmarshal([]byte(message), &receivedMsg); err != nil {
			panic(err)
		}

		// Parse Go map into entity.Orderbook
		messageType, ticker, orderBook := parseData(receivedMsg, exchange.maxDepth)
		// set exchange.OrderBooks based on snapshot or update messageType
		exchange.UpdateOrderBooks(messageType, ticker, orderBook)
		// fmt.Println("new message: ", messageType, ticker, orderBook)
		// fmt.Println("updated orderBooks: ", exchange.orderBooks["btcusdt"], exchange.orderBooks["ethusdt"])
	}
}

func (exchange *Exchange) UpdateOrderBooks(messageType string, ticker string, orderBook entity.Orderbook) {
	if messageType == "snapshot" {
		orderBook.Ts_updated = time.Now().UTC()
		exchange.orderBooks[ticker] = &orderBook
		// send ordeBbook to aggregator channel
	} else if messageType == "update" {
		// insert update Bids and Asks into orderBooks, trim to maxDepth
		if exchange.orderBooks[ticker] != nil {

			entity.InsertBids(&exchange.orderBooks[ticker].Bids, orderBook.Bids)
			exchange.orderBooks[ticker].Bids = exchange.orderBooks[ticker].Bids[:int(math.Min(float64(len(exchange.orderBooks[ticker].Bids)), float64(exchange.maxDepth)))]
			entity.InsertAsks(&exchange.orderBooks[ticker].Asks, orderBook.Asks)
			exchange.orderBooks[ticker].Asks = exchange.orderBooks[ticker].Asks[:int(math.Min(float64(len(exchange.orderBooks[ticker].Asks)), float64(exchange.maxDepth)))]

			// update orderBooks's Ts_original, Ts_received, Ts_updated from update message and current time respectively
			exchange.orderBooks[ticker].Ts_original = orderBook.Ts_original
			exchange.orderBooks[ticker].Ts_received = orderBook.Ts_received
			exchange.orderBooks[ticker].Ts_updated = time.Now().UTC()

			// send exchange.orderBooks[ticker] to aggregator channel
		}
	}
}

// func insertBids(existing *[]entity.PriceSize, update []entity.PriceSize) {

// 	for _, elem := range update {
// 		idx := sort.Search(len(*existing), func(i int) bool {
// 			return (*existing)[i].Price <= elem.Price
// 		})
// 		if idx == len(*existing) {
// 			if elem.Size != 0 {
// 				*existing = append((*existing), elem)
// 			}
// 		} else {

// 			if (*existing)[idx].Price == elem.Price {
// 				if elem.Size != 0 {
// 					*existing = append((*existing)[:idx], append([]entity.PriceSize{elem}, (*existing)[idx+1:]...)...)
// 				} else {
// 					*existing = append((*existing)[:idx], (*existing)[idx+1:]...)
// 				}
// 			} else if (*existing)[idx].Price < elem.Price {
// 				*existing = append((*existing)[:idx], append([]entity.PriceSize{elem}, (*existing)[idx:]...)...)
// 			}
// 		}

// 	}
// }

// func insertAsks(existing *[]entity.PriceSize, update []entity.PriceSize) {

// 	for _, elem := range update {
// 		idx := sort.Search(len(*existing), func(i int) bool {
// 			return (*existing)[i].Price >= elem.Price
// 		})
// 		if idx == len(*existing) {
// 			if elem.Size != 0 {
// 				*existing = append((*existing), elem)
// 			}
// 		} else {

// 			if (*existing)[idx].Price == elem.Price {
// 				if elem.Size != 0 {
// 					*existing = append((*existing)[:idx], append([]entity.PriceSize{elem}, (*existing)[idx+1:]...)...)
// 				} else {
// 					*existing = append((*existing)[:idx], (*existing)[idx+1:]...)
// 				}
// 			} else if (*existing)[idx].Price > elem.Price {
// 				*existing = append((*existing)[:idx], append([]entity.PriceSize{elem}, (*existing)[idx:]...)...)
// 			}
// 		}
// 	}
// }
