# Cchecker
# A simple content checker written in Go.

Cchecker is a command-line tool that checks for specific content in files or directories. It can be used to ensure that certain keywords or phrases are present or absent in your codebase or documents.

## Features
- Check for the presence of specific keywords in files.
- Recursively check directories.
- Supports wildcard ("*") or specific file extensions (e.g., ".go", ".txt").
- Outputs results to the console.
- Easy to use with command-line arguments.

## Installation
To install Cchecker, you need to have Go installed on your machine. Then, you can clone the repository and build the tool:
```bash
git clone git@github.com:MonkyMars/ccheck.git
cd ccheck
go build -o cchecker .
```

The cchecker binary will be created in the current directory.

## Usage
You can run Cchecker from the command line with the following syntax:
```bash
./cchecker <pattern> <root> <file_extension>
```

- `<pattern>`: The keyword or phrase you want to check for.
- `<root>`: The root directory to start the check.
- `<file_extension>`: The file extension to filter files (use "*" for all files).

Example:
```bash
./cchecker "TODO" ~/Coding ".go"
```

Output:
```
/home/monky/Coding/go/check/main.go:68: 	/// E.g., Pattern: TODO, root: /home/monky/go, ext: .go
```

## Contributing
Contributions are welcome! If you find a bug or have a feature request, please open an issue. If you'd like to contribute code, feel free to fork the repository and submit a pull request.

## License
This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
