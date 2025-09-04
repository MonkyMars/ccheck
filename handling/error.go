package handling

import "fmt"

func PrintError(msg string, expected string) string {
	return fmt.Sprintf("Error: %s, expected: %s", msg, expected)
}
