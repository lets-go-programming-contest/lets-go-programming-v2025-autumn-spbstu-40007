package die

import (
	"fmt"
	"os"
)

func Die(args ...any) {
	fmt.Fprintln(os.Stderr, args...)

	os.Exit(0)
}
