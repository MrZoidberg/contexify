# Contexify

## Overview

Contexify is a lightweight, fast, and flexible command-line tool designed to aggregate files in a folder structure and prepare a context file from all files within a specified folder subtree. It is useful for generating a single file containing all the contents of a project as a context for LLMs.

## Features

- **File Filtering**: Include or exclude files based on patterns and `.gitignore` rules.
- **Folder Tree Representation**: Generates a structured tree view of the file system.
- **Token Estimation**: Estimates token count for the final context.

## Installation

### Using Go

You can install Contexify directly from the source:

```sh
# Clone the repository
go install github.com/MrZoidberg/contexify
```

### Manual Build

Alternatively, you can manually download the release binaries from the [GitHub Releases](https://github.com/MrZoidberg/contexify/releases)

## Usage

Run the tool with the following options:

```sh
contexify -i <input-directory> -o <output-file> [options]
```

### Command-line Options

| Flag | Description | Default |
|------|------------|---------|
| `-i, --input` | Input folder path | `.` |
| `-o, --output` | Output file path | `context.txt` |
| `--include` | File patterns to include (semicolon-separated) | `` |
| `--exclude` | File patterns to exclude (semicolon-separated) | `LICENSE;CHANGELOG.md` |
| `--disable-gitignore` | Ignore `.gitignore` rules | `false` |
| `--disable-folder-tree` | Exclude folder tree from output | `false` |
| `--non-recursive` | Do not scan subdirectories | `false` |
| `--delimiter` | Custom delimiter between files | `---` |
| `--tokenizer.skip` | Skip token count estimation | `false` |
| `-v, --verbose` | Enable verbose logging | `false` |

## Example

```sh
contexify -i ./data -o output.txt --include "*.txt;*.md" --delimiter "\n---\n"
```

This command processes `data/`, extracts `.txt` and `.md` files, and writes the output to `output.txt` with a `---` delimiter.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss improvements.

## License

This project is licensed under the MIT License.

## Author

Maintained by **MrZoidberg**.

---
For more information, visit [GitHub Repository](https://github.com/MrZoidberg/contexify).
