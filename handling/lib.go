package handling

import (
	"fmt"
	"os"
	"path/filepath"
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

func HandleOutputFileArg(args []string) (*os.File, error) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-o") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 || parts[1] == "" {
				return nil, fmt.Errorf("invalid output file argument, expected -o=output.txt, got %s", arg)
			}

			// Validate the file path
			outputFilePath := parts[1]

			// Ensure the directory exists (if not, return an error)
			dir := filepath.Dir(outputFilePath)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return nil, fmt.Errorf("directory %s does not exist", dir)
			}

			// Try to create the output file
			outputFile, err := os.Create(outputFilePath)
			if err != nil {
				return nil, fmt.Errorf("unable to create output file %s: %v", outputFilePath, err)
			}

			// Return the created file pointer
			return outputFile, nil
		}
	}
	// Return nil and no error if no output argument is found
	return nil, nil
}
