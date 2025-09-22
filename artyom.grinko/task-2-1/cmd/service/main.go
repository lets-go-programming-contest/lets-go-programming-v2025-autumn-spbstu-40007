package main

//nolint:all
import (
	"fmt"
	"strconv"
	"strings"

	"task-2-1/internal/scanner"
)

const (
	minLowerBoundary = 15
	maxUpperBoundary = 30
)

func main() {
	scanner := scanner.NewScanner()
	n, _ := strconv.Atoi(scanner.Read())

	for range n {
		//nolint:all
		k, _ := strconv.Atoi(scanner.Read())

		lowerBoundary, upperBoundary := minLowerBoundary, maxUpperBoundary
		//nolint:all
		for j := range k {
			preferences := scanner.Read()
			mode, temperature := func() (string, int) {
				temporary := strings.Fields(preferences)
				temperature, _ := strconv.Atoi(temporary[1])

				return temporary[0], temperature
			}()

			if mode == ">=" {
				lowerBoundary = max(lowerBoundary, temperature)
			} else {
				upperBoundary = min(upperBoundary, temperature)
			}

			if lowerBoundary > upperBoundary {
				for range k - j {
					fmt.Println(-1)
				}

				scanner.SkipNLines(k - j - 1)

				break
			} else {
				fmt.Println(lowerBoundary)
			}
		}
	}
}
