// package main ...
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/codeYann/go-pagination/internal/requests"
)

type data struct {
	Amount float32 `json:"amount"`
	Date   int64   `json:"date"`
	Price  float64 `json:"price"`
	Tid    int64   `json:"tid"`
	Type   string  `json:"type"`
}

func main() {
	runtime.GOMAXPROCS(2)

	request := requests.HTTPRequest{
		BaseURL:    "https://www.mercadobitcoin.net/api/BTC/",
		Method:     "GET",
		MaxTimeout: 2 * time.Second,
	}

	url := "trades"

	ctx := context.TODO()

	res, err := request.MakeRequest(ctx, url)
	if err != nil {
		log.Println(err)
		return
	}

	var bitcoinsData []data

	err = json.Unmarshal(res, &bitcoinsData)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(bitcoinsData)
}
