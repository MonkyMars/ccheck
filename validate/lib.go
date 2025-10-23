package validate

import (
	"path/filepath"
	"strings"
)

func IsValidDir(path string, blacklisted_dirs map[string]bool) bool {
	_, is_blacklisted_dir := blacklisted_dirs[filepath.Base(path)]
	is_hiddendir := strings.HasPrefix(filepath.Base(path), ".")
	return !is_blacklisted_dir && !is_hiddendir
}

func IsValidExtension(file, ext string) bool {
	if ext == "*" {
		return true
	}
	return filepath.Ext(file) == ext
}
