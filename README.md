# Cchecker
# A lightweight content checker written in Go.

Cchecker is a command-line tool that checks for specific content in files or directories. It can be used to ensure that certain keywords or phrases are present or absent in your codebase or documents.

## Features
- Extremely fast and lightweight.
- Check for the presence of specific keywords in files.
- Supports regex patterns (prefix pattern with `re:`).
- Case-sensitive and case-insensitive search options.
- Recursively check directories.
- Supports wildcard ("*") or specific file extensions (e.g., ".go", ".txt").
- Easy to use with command-line arguments.

## Installation
To install Cchecker, you need to have Go installed on your machine. Then, you can clone the repository and build the tool:
```bash
git clone git@github.com:MonkyMars/ccheck.git
cd ccheck
go build -o ccheck .
```

The cchecker binary will be created in the current directory.

## Usage
You can run Cchecker from the command line with the following syntax:
```bash
./ccheck <pattern> <root> <file_extension> <flags>
```

- `<pattern>`: The keyword or phrase you want to check for, can be multiple words, separated by a comma. Prefix with `re:` to treat it as a regex pattern.
- `<root>`: The root directory to start the check.
- `<file_extension>`: The file extension to filter files (use "*" for all files, separate by a comma for multiple file extentions).
- `<flags>`: Optional flags for additional functionality. Currently only supports -i for case-insensitive search and -o=<file_path> for saving the result to a file.

Example:
```bash
./ccheck "TODO, re:^func" ~/Coding ".go, .rs" -i -o=output.txt
```
or
```bash
./ccheck "re:TODO" ~/Coding ".go" -i
```

Output:
```
Cchecker x.y.z
/home/monky/Coding/go/check/main.go:68: 	/// E.g., Pattern: TODO, root: /home/monky/go, ext: .go
```

> Note:
- Be careful with setting the root directory to something too broad. For example, using ~/ as the root, will check all your package files too. node_modules for javascript and target for rust are already filtered out by default. You can modify the code to add more directories to ignore if needed.

## Contributing
Contributions are welcome! If you find a bug or have a feature request, please open an issue. If you'd like to contribute code, feel free to fork the repository and submit a pull request.

## License
This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
