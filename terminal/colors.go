package terminal

import (
	"fmt"
	"regexp"
	"strings"
)

// ANSI color codes
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[90m"
var White = "\033[97m"

// ColorTheme defines colors for different parts of the output
type ColorTheme struct {
	FilePath   string
	LineNumber string
	MatchText  string
	NormalText string
}

// DefaultTheme provides the default color scheme
var DefaultTheme = ColorTheme{
	FilePath:   Cyan,
	LineNumber: Yellow,
	MatchText:  Red,
	NormalText: Gray,
}

// Colorize wraps text with ANSI color codes
func Colorize(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// ColorizeMatch highlights all pattern matches in a line of text
func ColorizeMatch(line string, pattern *regexp.Regexp, theme ColorTheme) string {
	if pattern == nil {
		return Colorize(line, theme.NormalText)
	}

	// Find all matches with their positions
	matches := pattern.FindAllStringIndex(line, -1)
	if len(matches) == 0 {
		return Colorize(line, theme.NormalText)
	}

	var result strings.Builder
	lastEnd := 0

	for _, match := range matches {
		start, end := match[0], match[1]

		// Add text before the match
		if start > lastEnd {
			result.WriteString(Colorize(line[lastEnd:start], theme.NormalText))
		}

		// Add the highlighted match
		result.WriteString(Colorize(line[start:end], theme.MatchText))

		lastEnd = end
	}

	// Add remaining text after the last match
	if lastEnd < len(line) {
		result.WriteString(Colorize(line[lastEnd:], theme.NormalText))
	}

	return result.String()
}

// FormatMatchLine creates a formatted output line for a pattern match
func FormatMatchLine(filePath string, lineNumber int, lineText string, pattern *regexp.Regexp) string {
	pathText := Colorize(filePath, DefaultTheme.FilePath)
	lineNumText := Colorize(fmt.Sprintf("%d", lineNumber), DefaultTheme.LineNumber)
	coloredLineText := ColorizeMatch(lineText, pattern, DefaultTheme)

	return fmt.Sprintf("%s:%s: %s", pathText, lineNumText, coloredLineText)
}
