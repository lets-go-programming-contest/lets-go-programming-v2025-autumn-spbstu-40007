package scanner

import (
	"bufio"
	"os"

	"task-2-1/internal/die"
)

type Scanner struct {
	*bufio.Scanner
}

func NewScanner() *Scanner {
	return &Scanner{bufio.NewScanner(os.Stdin)}
}

func (scanner *Scanner) Read() string {
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		die.Die(err)
	}

	return scanner.Text()
}

func (scanner *Scanner) SkipNLines(n int) {
	for range n {
		scanner.Read()
	}
}
