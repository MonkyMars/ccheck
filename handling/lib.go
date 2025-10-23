package handling

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// HandleCaseSensitivityArg checks for the presence of the -i flag in the command-line arguments
func HandleCaseSensitivityArg() bool {
	return !slices.Contains(os.Args[4:], "-i")
}

// HandleOutputFileArg processes the -o=output.txt argument to create and return an output file
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

// HandleHelpArg checks for -h or --help flags and displays usage information
func HandleHelpArg(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			fmt.Println("Usage: go run main.go <pattern|re:regex> <root_dir> <ext> <flags>")
			fmt.Println("Arguments:")
			fmt.Println("  <pattern|re:regex> : Comma-separated list of literal patterns or regex (prefix with re:)")
			fmt.Println("  <root_dir>        : Root directory to search")
			fmt.Println("  <ext>             : Comma-separated list of file extensions to include")
			fmt.Println("Flags:")
			fmt.Println("  -i                : Case insensitive search")
			fmt.Println("  -o=output.txt    : Specify output file")
			fmt.Println("  -h, --help       : Show this help message")
			os.Exit(0)
		}
	}
	return false
}
