package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().UnixNano())
	randomNum := random(5, 10)
	wordPtr := flag.String("url", "foo", "a string")

	flag.Parse()

	fmt.Printf("Random Num: %d\n", randomNum)
	time.Sleep(time.Duration(randomNum) * time.Second)

	url := "http://" + *wordPtr
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		code := response.StatusCode

		fmt.Printf("Code: %d\n", code)
	}
}
