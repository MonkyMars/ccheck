package handling

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ParseArgs() (pattern *regexp.Regexp, root string, extList []string, outputFile *os.File) {
	if len(os.Args) < 4 {
		fmt.Println(PrintError("not enough arguments", "at least 3 arguments required"))
		fmt.Println("Usage: go run main.go <pattern|re:regex> <root_dir> <ext> <flags>")
		os.Exit(1)
	}

	case_sensitive := HandleCaseSensitivityArg()
	outputFile, err := HandleOutputFileArg(os.Args[4:])
	if err != nil {
		fmt.Println(PrintError(err.Error(), "valid output file argument"))
		os.Exit(1)
	}
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

	root = os.Args[2]
	if root == "." {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(PrintError(err.Error(), "get current working directory"))
			os.Exit(1)
		}
		root = cwd
	}

	ext := os.Args[3]
	extList = strings.Split(ext, ",")
	for i, e := range extList {
		extList[i] = strings.TrimSpace(e)
		if !strings.HasPrefix(extList[i], ".") {
			extList[i] = "." + extList[i]
		}
	}

	/// E.g., Pattern: TODO, root: /home/monky/go, ext: .go
	return s, root, extList, outputFile
}
