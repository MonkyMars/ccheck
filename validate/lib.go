package validate

import (
	"path/filepath"
	"slices"
	"strings"
)

func Is_valid_dir(path string, blacklisted_dirs []string) bool {
	is_blacklisted_dir := slices.Contains(blacklisted_dirs, filepath.Base(path))
	is_hiddendir := strings.HasPrefix(filepath.Base(path), ".")
	return !is_blacklisted_dir && !is_hiddendir
}

func Is_valid_ext(file, ext string) bool {
	if ext == "*" {
		return true
	}
	return filepath.Ext(file) == ext
}
