package main

import (
	"bufio"
	"ccheck/handling"
	"ccheck/validate"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

var blacklistedDirs = []string{"node_modules", "target"}
var validFlags = []string{"-i", "-o"}

func main() {
	fmt.Println("Cchecker 1.2.1")
	pattern, root, ext, outputFile := handling.ParseArgs()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(handling.PrintError(err.Error(), "file should be accessible"))
			return nil
		}

		if d.IsDir() {
			if !validate.Is_valid_dir(path, blacklistedDirs) {
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
		file, err := handling.OpenFile(path)
		if err != nil {
			fmt.Println(handling.PrintError(err.Error(), "file should be accessible"))
			return nil
		}

		scanner := bufio.NewScanner(file)
		lineNum := 1

		for scanner.Scan() {
			line := scanner.Text()
			if pattern.MatchString(line) {
				message := fmt.Sprintf("%s:%d: %s\n", path, lineNum, line)
				handling.OutputToFile(outputFile, message)
			}
			lineNum++
		}

		err = file.Close()
		if err != nil {
			fmt.Println(handling.PrintError(err.Error(), "file should be closed"))
			return nil
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(handling.PrintError(err.Error(), "error reading file"))
		}

		return nil
	})

	if err != nil {
		fmt.Println(handling.PrintError(err.Error(), "error walking the directory"))
	}
}
