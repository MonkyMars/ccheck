package main

import (
	"ccheck/validate"
	"testing"
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
		result := validate.Is_valid_ext(tt.file, tt.ext)
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
		result := validate.Is_valid_dir(tt.dir, blacklisted_dirs)
		if result != tt.expected {
			t.Errorf("isValidDir(%q) = %v; want %v", tt.dir, result, tt.expected)
		}
	}
}
