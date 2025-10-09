package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MinTemp = 15
	MaxTemp = 30
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		os.Exit(1)
	}
	N, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || N < 1 || N > 1000 {
		os.Exit(1)
	}

	for dept := 0; dept < N; dept++ {
		if !scanner.Scan() {
			os.Exit(1)
		}
		K, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || K < 1 || K > 1000 {
			os.Exit(1)
		}

		minAllowed := MinTemp
		maxAllowed := MaxTemp
		failed := false

		for emp := 0; emp < K && scanner.Scan(); emp++ {
			if failed {
				continue
			}

			line := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(line)
			if len(parts) != 2 {
				os.Exit(1)
			}

			op := parts[0]
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				os.Exit(1)
			}

			switch op {
			case ">=":
				if val > minAllowed {
					minAllowed = val
				}
			case "<=":
				if val < maxAllowed {
					maxAllowed = val
				}
			default:
				os.Exit(1)
			}

			if minAllowed <= maxAllowed {
				fmt.Println(minAllowed)
			} else {
				fmt.Println(-1)
				failed = true
			}
		}
	}
}
