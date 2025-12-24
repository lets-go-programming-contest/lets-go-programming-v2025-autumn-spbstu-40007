package main

import (
	"fmt"
	"os"

	"github.com/ami0-0/config"
)

func main() {
	params, err := config.Fetch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	fmt.Print(params.String())
}
