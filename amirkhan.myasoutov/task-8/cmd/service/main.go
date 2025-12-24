package main

import (
	"fmt"
	"log"

	"github.com/ami0-0/config"
)

func main() {
	params, err := config.Fetch()
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
		return
	}

	fmt.Println(params)
}
