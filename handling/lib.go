package handling

import (
	"fmt"
	"os"
	"strings"
)

func HandleCaseSensitivityArg() bool {
	case_sensitive := true
	for _, arg := range os.Args[4:] {
		if arg == "-i" {
			case_sensitive = false
		}
	}
	return case_sensitive
}

func HandleOutputFileArg() *os.File {
	for _, arg := range os.Args[4:] {
		if strings.HasPrefix(arg, "-o") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 || parts[1] == "" {
				fmt.Println(PrintError("invalid output file argument", "-o=output.txt"))
				os.Exit(1)
			}
			outputFile, err := os.Create(parts[1])
			if err != nil {
				fmt.Println(PrintError(err.Error(), "output file should be creatable"))
				os.Exit(1)
			}
			return outputFile
		}
	}
	return nil
}
