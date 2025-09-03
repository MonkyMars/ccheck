package main

import (
	"bufio"
	"ccheck/validate"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

var blacklisted_dirs = []string{"node_modules", "target"}

func main() {
	fmt.Println("Cchecker 1.2.1")
	pattern, root, ext := parseArgs()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(printError(err.Error(), "file should be accessible"))
			return nil
		}

		if d.IsDir() {
			if !validate.Is_valid_dir(path, blacklisted_dirs) {
				return filepath.SkipDir
			}
			return nil
		}

		if !validate.Is_valid_ext(d.Name(), ext) {
			return nil
		}

		// Ensure path is inside root
		rel, err := filepath.Rel(root, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			return nil // skip paths outside root
		}

		// #nosec G304: Path is validated to be inside the root directory
		file, err := openFile(path)
		if err != nil {
			fmt.Println(printError(err.Error(), "file should be accessible"))
			return nil
		}

		scanner := bufio.NewScanner(file)
		lineNum := 1

		for scanner.Scan() {
			line := scanner.Text()
			if pattern.MatchString(line) {
				fmt.Printf("%s:%d: %s\n", path, lineNum, line)
			}
			lineNum++
		}

		err = file.Close()
		if err != nil {
			fmt.Println(printError(err.Error(), "file should be closed"))
			return nil
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(printError(err.Error(), "error reading file"))
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}

func parseArgs() (pattern *regexp.Regexp, root string, ext string) {
	if len(os.Args) < 4 {
		fmt.Println(printError("not enough arguments", "at least 3 arguments required"))
		fmt.Println("Usage: go run main.go <pattern|re:regex> <root_dir> <ext> <flags>")
		os.Exit(1)
	}

	case_sensitive := true
	for _, arg := range os.Args[4:] {
		if arg == "-i" {
			case_sensitive = false
		} else {
			fmt.Println(printError("unknown flag "+arg, "valid flags are -i"))
			os.Exit(1)
		}
	}
	var s *regexp.Regexp
	patternArg := os.Args[1]

	if strings.HasPrefix(patternArg, "re:") {
		regexPattern := patternArg[3:] // remove "re:"
		if !case_sensitive {
			regexPattern = "(?i)" + regexPattern
		}
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			fmt.Println(printError(err.Error(), "valid regex"))
			os.Exit(1)
		}
		s = re
	} else {
		// Literal search
		literal := regexp.QuoteMeta(patternArg)
		if !case_sensitive {
			literal = "(?i)" + literal
		}
		s = regexp.MustCompile(literal)
	}

	/// E.g., Pattern: TODO, root: /home/monky/go, ext: .go
	return s, os.Args[2], os.Args[3]
}

func printError(msg string, expected string) string {
	return fmt.Sprintf("Error: %s, expected: %s", msg, expected)
}

func isBinaryFile(file *os.File) bool {
	buf := make([]byte, 8000)
	n, err := file.Read(buf)
	if err != nil {
		return true // unreadable → assume binary
	}
	if n == 0 {
		return false // empty file → treat as text
	}

	// Look for null bytes
	for _, b := range buf[:n] {
		if b == 0 {
			return true
		}
	}

	// Validate UTF-8
	if !utf8.Valid(buf[:n]) {
		return true
	}

	return false
}

func openFile(path string) (*os.File, error) {
	// #nosec G304: Path is validated to be inside the root directory
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if isBinaryFile(file) {
		_ = file.Close()
		return nil, fmt.Errorf("binary file")
	}

	// Reset file cursor after binary check
	_, err = file.Seek(0, 0)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	return file, nil
}
