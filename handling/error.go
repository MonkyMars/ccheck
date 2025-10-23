package handling

import (
	"errors"
	"fmt"
)

var ErrorBinaryFile = errors.New("binary file")

func PrintError(msg string, expected string) string {
	return fmt.Sprintf("Error: %s, expected: %s", msg, expected)
}
