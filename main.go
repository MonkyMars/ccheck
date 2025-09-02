package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
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
			if !is_valid_dir(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if !is_valid_ext(d.Name(), ext) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Println(print_error(err.Error(), "file should be accessible"))
			return nil
			// TODO: handle error properly
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 1

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, pattern) {
				fmt.Printf("%s:%d: %s\n", path, lineNum, line)
			}
			lineNum++
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}

func parse_args() (pattern string, root string, ext string) {
	if len(os.Args) < 4 {
		fmt.Println(print_error("not enough arguments", "at least 3 arguments required"))
		os.Exit(1)
	}
	/// Pattern, root directory file extension
	/// E.g., Pattern: TODO, root: /home/monky/go, ext: .go
	return os.Args[1], os.Args[2], os.Args[3]
}

func print_error(msg string, expected string) string {
	return fmt.Sprintf("Error: %s, expected: %s", msg, expected)
}

func is_valid_dir(path string) bool {
	is_blacklisted_dir := slices.Contains(blacklisted_dirs, filepath.Base(path))
	is_hiddendir := strings.HasPrefix(filepath.Base(path), ".")
	return !is_blacklisted_dir && !is_hiddendir
}

func is_valid_ext(file, ext string) bool {
	if ext == "*" {
		return true
	}
	return filepath.Ext(file) == ext
}
