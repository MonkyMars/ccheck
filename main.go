package main

import (
	"bufio"
	"ccheck/handling"
	"ccheck/validate"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	blacklistedDirs = []string{"node_modules", "target", ".git"}
	results         = make(chan string)
	resultsLen      = 0
)

func main() {
	fmt.Println("ccheck 2.1.0")
	patterns, root, extList, outputFile := handling.ParseArgs()

	// Start timer
	startTime := time.Now()

	var wg sync.WaitGroup
	var outputWg sync.WaitGroup

	// Start a goroutine to handle output
	outputWg.Go(func() {
		handling.OutputToFile(outputFile, results)
	})

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

		matches := false
		for _, e := range extList {
			if validate.Is_valid_ext(d.Name(), e) {
				matches = true
				break
			}
		}
		if !matches {
			return nil
		}

		// Avoid processing the output file for an infinite loop
		if outputFile != nil && path == outputFile.Name() {
			return nil
		}

		// Ensure path is inside root
		rel, err := filepath.Rel(root, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			return nil // skip paths outside root
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			// #nosec G304: Path is validated to be inside the root directory
			file, err := handling.OpenFile(path)
			if err != nil || file == nil {
				if err.Error() != "binary file" {
					fmt.Println(handling.PrintError(err.Error(), "file should be accessible"))
				}
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNum := 1

			for scanner.Scan() {
				line := scanner.Text()
				for _, p := range patterns {
					if p.MatchString(line) {
						message := fmt.Sprintf("%s:%d: %s\n", path, lineNum, line)
						results <- message
						resultsLen++
						break // No need to check other patterns if one matches
					}
				}
				lineNum++
			}

			if err := scanner.Err(); err != nil {
				fmt.Println(handling.PrintError(err.Error(), "error reading file"))
			}
		}(path)
		return nil
	})

	wg.Wait()

	close(results)
	outputWg.Wait()

	elapsed := time.Since(startTime)
	fmt.Printf("Found %d results in %s\n", resultsLen, elapsed)

	if err != nil {
		fmt.Println(handling.PrintError(err.Error(), "error walking the directory"))
	}
}
