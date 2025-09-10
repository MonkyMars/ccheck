package handling

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func HandleCaseSensitivityArg() bool {
	return !slices.Contains(os.Args[4:], "-i")
}

func HandleOutputFileArg(args []string) (*os.File, error) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-o") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 || parts[1] == "" {
				return nil, fmt.Errorf("invalid output file argument, expected -o=output.txt, got %s", arg)
			}

			outputFilePath := parts[1]
			cleanPath := filepath.Clean(outputFilePath)

			// Reject if the cleaned path is a directory traversal
			if cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(os.PathSeparator)) {
				return nil, fmt.Errorf("invalid output file path: directory traversal detected")
			}

			// Get absolute path
			absPath := outputFilePath
			if !filepath.IsAbs(outputFilePath) {
				cwd, _ := os.Getwd()
				absPath = filepath.Join(cwd, outputFilePath)
			}
			absPath = filepath.Clean(absPath)

			// Ensure the directory exists
			dir := filepath.Dir(absPath)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return nil, fmt.Errorf("directory %s does not exist", dir)
			}

			outputFile, err := os.Create(absPath)
			if err != nil {
				return nil, fmt.Errorf("unable to create output file %s: %v", absPath, err)
			}
			return outputFile, nil
		}
	}
	return nil, nil // No output file specified
}
