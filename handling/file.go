package handling

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func IsBinaryFile(file *os.File) bool {
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

func OpenFile(path string) (*os.File, error) {
	// #nosec G304: Path is validated to be inside the root directory
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if IsBinaryFile(file) {
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

func OutputToFile(outputFile *os.File, message string) {
	if outputFile != nil {
		_, err := outputFile.WriteString(message + "\n")
		if err != nil {
			PrintError(err.Error(), "output file should be writable")
			os.Exit(1)
		}
	} else {
		fmt.Println(message)
	}
}
