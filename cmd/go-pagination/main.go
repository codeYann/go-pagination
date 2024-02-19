// package main ...
package main

import (
	"fmt"
	"runtime"

	"github.com/codeYann/go-pagination/internal/requests"
)

func main() {
	runtime.GOMAXPROCS(1)

	request := requests.HTTPRequest{
		BaseURL: "https://jsonplaceholder.typicode.com/",
		Method:  "GET",
	}

	url := "posts"

	data, err := request.Request(url)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}
