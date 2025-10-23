package main

import (
	"testing"

	"github.com/MonkyMars/ccheck/validate"
)

func TestIsValidExt(t *testing.T) {
	tests := []struct {
		file     string
		ext      string
		expected bool
	}{
		{"file.go", ".go", true},
		{"file.txt", ".go", false},
		{"file.go", "*", true},
	}

	for _, tt := range tests {
		result := validate.IsValidExtension(tt.file, tt.ext)
		if result != tt.expected {
			t.Errorf("isValidExt(%q, %q) = %v; want %v", tt.file, tt.ext, result, tt.expected)
		}
	}
}

func TestIsValidDir(t *testing.T) {
	tests := []struct {
		dir      string
		expected bool
	}{
		{"node_modules", false},
		{"target", false},
		{"src", true},
		{"lib", true},
	}

	for _, tt := range tests {
		result := validate.IsValidDir(tt.dir, blacklistedDirs)
		if result != tt.expected {
			t.Errorf("isValidDir(%q) = %v; want %v", tt.dir, result, tt.expected)
		}
	}
}
