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
)

var blacklisted_dirs = []string{"node_modules", "target"}

func main() {
	pattern, root, ext := parse_args()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(print_error(err.Error(), "file should be accessible"))
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
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(print_error(err.Error(), "file should be accessible"))
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
			fmt.Println(print_error(err.Error(), "file should be closed"))
			return nil
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(print_error(err.Error(), "error reading file"))
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}

func parse_args() (pattern *regexp.Regexp, root string, ext string) {
	if len(os.Args) < 4 {
		fmt.Println(print_error("not enough arguments", "at least 3 arguments required"))
		fmt.Println("Usage: go run main.go <pattern|re:regex> <root_dir> <ext> <flags>")
		os.Exit(1)
	}

	case_sensitive := true
	for _, arg := range os.Args[4:] {
		if arg == "-i" {
			case_sensitive = false
		} else {
			fmt.Println(print_error("unknown flag "+arg, "valid flags are -i"))
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
			fmt.Println(print_error(err.Error(), "valid regex"))
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

func print_error(msg string, expected string) string {
	return fmt.Sprintf("Error: %s, expected: %s", msg, expected)
}
