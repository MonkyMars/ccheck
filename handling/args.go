package handling

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ParseArgs() (pattern *regexp.Regexp, root string, ext string, outputFile *os.File) {
	if len(os.Args) < 4 {
		fmt.Println(PrintError("not enough arguments", "at least 3 arguments required"))
		fmt.Println("Usage: go run main.go <pattern|re:regex> <root_dir> <ext> <flags>")
		os.Exit(1)
	}

	case_sensitive := HandleCaseSensitivityArg()
	outputFile = HandleOutputFileArg()
	var s *regexp.Regexp
	patternArg := os.Args[1]

	if strings.HasPrefix(patternArg, "re:") {
		regexPattern := patternArg[3:] // remove "re:"
		if !case_sensitive {
			regexPattern = "(?i)" + regexPattern
		}
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			fmt.Println(PrintError(err.Error(), "valid regex"))
			os.Exit(1)
		}
		s = re
	} else {
		// Literal search
		literal := regexp.QuoteMeta(patternArg)
		if !case_sensitive {
			literal = "(?i)" + literal
		}
		s = regexp.MustCompile(literal)
	}

	/// E.g., Pattern: TODO, root: /home/monky/go, ext: .go
	return s, os.Args[2], os.Args[3], outputFile
}
