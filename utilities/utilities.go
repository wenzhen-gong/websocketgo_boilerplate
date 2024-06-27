package utilities

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateKraken(data map[string]interface{}) bool {

	validate := validator.New()

	rules := map[string]interface{}{
		"channel": "required,eq=book",
		"type":    "required,eq=update|eq=snapshot",
		"data":    "required",
	}

	data0rules := map[string]interface{}{
		"symbol":    "required,eq=BTC/USDT|eq=ETH/USDT",
		"checksum":  "required,numeric",
		"timestamp": "required",
		"bids":      "required",
		"asks":      "required",
	}

	bidsasksrules := map[string]interface{}{
		"price": "gte=0",
		"qty":   "gte=0",
	}

	errs1 := validate.ValidateMap(data, rules)
	if len(errs1) > 0 {
		fmt.Println("errs1")
		return false
	}
	errs2 := validate.ValidateMap(data["data"].([]interface{})[0].(map[string]interface{}), data0rules)
	if len(errs2) > 0 {
		fmt.Println("errs2")
		return false
	}

	for _, price_size := range data["data"].([]interface{})[0].(map[string]interface{})["asks"].([]interface{}) {

		err := validate.ValidateMap(price_size.(map[string]interface{}), bidsasksrules)
		if len(err) != 0 {
			fmt.Println("asks wrong")

			return false
		}
	}
	for _, price_size := range data["data"].([]interface{})[0].(map[string]interface{})["bids"].([]interface{}) {

		err := validate.ValidateMap(price_size.(map[string]interface{}), bidsasksrules)
		if len(err) != 0 {
			fmt.Println("bids wrong")

			return false
		}
	}
	return true
}
