package kraken

import (
	"math"
	"os"
	"time"
	"wz/entity"

	"github.com/mitchellh/mapstructure"
)

// should serve as a template for other exchanges, other parse functions need to return the same data format
func ParseKrakenData(receivedMsg map[string]interface{}, maxDepth int) (string, string, entity.Orderbook) {

	//<--- SHOULD DELETE output data receiving info
	fo, err := os.Create("from_kraken_received.log")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	//--->

	var msg krakenMsg

	errmsg := mapstructure.Decode(receivedMsg, &msg)
	if errmsg != nil {
		panic(errmsg)
	}

	layout := "2006-01-02T15:04:05.999999Z"

	if msg.Channel == "book" && msg.Type != "" {

		switch msg.Type {
		case "snapshot":

			//<--- SHOULD DELETE output data receiving info
			if _, err := fo.Write([]byte("Snapshot received ts: " + time.Now().UTC().String() + "\n")); err != nil {
				panic(err)
			}
			//--->

			// return the trimed Bids and Asks
			return "snapshot", entity.TickerMap[msg.Data[0].Symbol], entity.Orderbook{Bids: msg.Data[0].Bids[:int(math.Min(float64(len(msg.Data[0].Bids)), float64(maxDepth)))], Asks: msg.Data[0].Asks[:int(math.Min(float64(len(msg.Data[0].Asks)), float64(maxDepth)))], Ts_original: time.Now().UTC(), Ts_received: time.Now().UTC(), Ts_updated: time.Time{}}
		case "update":

			parsedTime, err := time.Parse(layout, msg.Data[0].Timestamp)
			if err != nil {
				panic(err)
			}

			//<--- SHOULD DELETE output data receiving info
			if _, err := fo.Write([]byte("Update received ts: " + parsedTime.String() + "\n")); err != nil {
				panic(err)
			}
			//--->

			return "update", entity.TickerMap[msg.Data[0].Symbol], entity.Orderbook{Bids: msg.Data[0].Bids, Asks: msg.Data[0].Asks, Ts_original: parsedTime, Ts_received: time.Now().UTC(), Ts_updated: time.Time{}}
		default:

			//<--- SHOULD DELETE output data receiving info
			if _, err := fo.Write([]byte("FAILED TO RECEIVE VALID MESSAGE\n")); err != nil {
				panic(err)
			}
			//--->

			return "", "", entity.Orderbook{}
		}
	}
	return "", "", entity.Orderbook{}
}

// helper structs
type krakenMsg struct {
	Channel string
	Type    string
	Data    []krakenMsgData
}

type krakenMsgData struct {
	Symbol    string
	Bids      []entity.PriceSize
	Asks      []entity.PriceSize
	Checksum  int
	Timestamp string
}

// func ValidateKraken(data map[string]interface{}) bool {

// 	validate := validator.New()

// 	rules := map[string]interface{}{
// 		"channel": "required,eq=book",
// 		"type":    "required,eq=update|eq=snapshot",
// 		"data":    "required",
// 	}

// 	data0rules := map[string]interface{}{
// 		"symbol":    "required,eq=BTC/USDT|eq=ETH/USDT",
// 		"checksum":  "required,numeric",
// 		"timestamp": "required",
// 		"bids":      "required",
// 		"asks":      "required",
// 	}

// 	bidsasksrules := map[string]interface{}{
// 		"price": "gte=0",
// 		"qty":   "gte=0",
// 	}

// 	errs1 := validate.ValidateMap(data, rules)
// 	if len(errs1) > 0 {
// 		// fmt.Println("errs1")
// 		return false
// 	}
// 	errs2 := validate.ValidateMap(data["data"].([]interface{})[0].(map[string]interface{}), data0rules)
// 	if len(errs2) > 0 {
// 		// fmt.Println("errs2")
// 		return false
// 	}

// 	for _, price_size := range data["data"].([]interface{})[0].(map[string]interface{})["asks"].([]interface{}) {

// 		err := validate.ValidateMap(price_size.(map[string]interface{}), bidsasksrules)
// 		if len(err) != 0 {
// 			// fmt.Println("asks wrong")

// 			return false
// 		}
// 	}
// 	for _, price_size := range data["data"].([]interface{})[0].(map[string]interface{})["bids"].([]interface{}) {

// 		err := validate.ValidateMap(price_size.(map[string]interface{}), bidsasksrules)
// 		if len(err) != 0 {
// 			// fmt.Println("bids wrong")

// 			return false
// 		}
// 	}
// 	return true
// }
