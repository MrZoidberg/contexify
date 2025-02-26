package app

import (
	"path/filepath"
	"sort"
	"strings"
)

// GenerateFileTree generates a string representation of a file tree
func GenerateFileTree(tree map[string][]string) (string, error) {
	var sb strings.Builder
	sb.WriteString("File Tree:\n")

	// Get folders and sort them
	folders := make([]string, 0, len(tree))
	for folder := range tree {
		folders = append(folders, folder)
	}
	sort.Strings(folders)

	// Process each folder
	for _, folder := range folders {
		sb.WriteString("└── " + folder + "\n")
		files := tree[folder]
		sort.Strings(files)

		// Process each file in the folder
		for i, file := range files {
			prefix := "    "
			if i == len(files)-1 {
				sb.WriteString(prefix + "└── " + filepath.Base(file) + "\n")
			} else {
				sb.WriteString(prefix + "├── " + filepath.Base(file) + "\n")
			}
		}
	}

	return sb.String(), nil
}
