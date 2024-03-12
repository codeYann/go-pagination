// package main ...
package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/codeYann/go-pagination/internal/requests"
)

func main() {
	runtime.GOMAXPROCS(2)

	request := requests.HTTPRequest{
		BaseURL:    "https://www.mercadobitcoin.net/api/BTC/",
		Method:     "GET",
		MaxTimeout: 300 * time.Millisecond,
	}

	url := "trades"

	ctx := context.TODO()

	data, err := request.MakeRequest(ctx, url)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(data))
}
