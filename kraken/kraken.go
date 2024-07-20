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
