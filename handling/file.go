package handling

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func IsBinaryFile(file *os.File) bool {
	// Check if the file has a file extension commonly associated with binary files
	ext := filepath.Ext(file.Name())
	if ext == "" || ext == ".exe" {
		return true
	}
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil || n == 0 {
		return false
	}

	// Quick null byte check
	for i := range n {
		if buf[i] == 0 {
			return true
		}
	}

	// UTF-8 validation only if needed
	return !utf8.Valid(buf[:n])
}

func OpenFile(path string) (*os.File, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	// Check if file is binary, if so, close and return error
	if IsBinaryFile(file) {
		_ = file.Close()
		return nil, ErrorBinaryFile
	}

	// Reset file cursor after binary check
	_, err = file.Seek(0, 0)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	return file, nil
}

func OutputToFile(outputFile *os.File, results chan string) {
	if outputFile != nil {
		for result := range results {
			_, err := outputFile.WriteString(result)
			if err != nil {
				fmt.Println(PrintError(err.Error(), "error writing to file"))
			}
		}
	} else {
		for result := range results {
			fmt.Print(result)
		}
	}
}
