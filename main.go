package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/MonkyMars/ccheck/handling"
	"github.com/MonkyMars/ccheck/terminal"
	"github.com/MonkyMars/ccheck/validate"
)

var (
	blacklistedDirs = map[string]bool{
		".git":         true,
		"node_modules": true,
		"target":       true,
	}
	results    = make(chan string)
	resultsLen = 0
)

func main() {
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
			if !validate.IsValidDir(path, blacklistedDirs) {
				return filepath.SkipDir
			}
			return nil
		}

		matches := false
		for e, _ := range extList {
			if validate.IsValidExtension(d.Name(), e) {
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

			file, err := handling.OpenFile(filepath.Clean(path))
			// Check for errors and nil file
			if err != nil || file == nil {
				if err.Error() != handling.ErrorBinaryFile.Error() {
					fmt.Println(handling.PrintError(err.Error(), "file should be accessible"))
				}
				return
			}
			defer file.Close()

			// Read file line by line
			scanner := bufio.NewScanner(file)
			lineNum := 1

			for scanner.Scan() {
				line := scanner.Text()
				for _, p := range patterns {
					if p.MatchString(line) {
						if outputFile != nil {
							message := fmt.Sprintf("%s:%d: %s\n", path, lineNum, line)
							results <- message
							resultsLen++
						} else {
							message := terminal.FormatMatchLine(path, lineNum, line, p)
							results <- message + "\n"
							resultsLen++
						}
						// Do not break for loop; another pattern might also match
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
	fmt.Println("ccheck 2.3.0")
	fmt.Printf("Found %d results in %s\n", resultsLen, elapsed)

	if err != nil {
		fmt.Println(handling.PrintError(err.Error(), "error walking the directory"))
	}
}
