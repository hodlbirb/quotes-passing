package main

import (
	"fmt"

	"github.com/FlashBoys/go-finance"
)

func main() {
	// 15-min delayed full quote for Apple.
	q, err := finance.GetQuote("AAPL")
	if err == nil {
		fmt.Println(q)
	}
}