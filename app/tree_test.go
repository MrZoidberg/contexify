package app

import (
	"strings"
	"testing"
)

func TestGenerateFileTree(t *testing.T) {
	tests := []struct {
		name     string
		tree     map[string][]string
		expected string
		wantErr  bool
	}{
		{
			name:     "Empty tree",
			tree:     map[string][]string{},
			expected: "File Tree:\n",
			wantErr:  false,
		},
		{
			name: "Single folder with no files",
			tree: map[string][]string{
				"folder1": {},
			},
			expected: "File Tree:\n└── folder1\n",
			wantErr:  false,
		},
		{
			name: "Single folder with one file",
			tree: map[string][]string{
				"folder1": {"file1.txt"},
			},
			expected: "File Tree:\n└── folder1\n    └── file1.txt\n",
			wantErr:  false,
		},
		{
			name: "Single folder with multiple files",
			tree: map[string][]string{
				"folder1": {"file1.txt", "file2.txt", "file3.txt"},
			},
			expected: "File Tree:\n└── folder1\n    ├── file1.txt\n    ├── file2.txt\n    └── file3.txt\n",
			wantErr:  false,
		},
		{
			name: "Multiple folders with multiple files",
			tree: map[string][]string{
				"folder1": {"file1.txt", "file2.txt"},
				"folder2": {"file3.txt", "file4.txt"},
				"folder3": {"file5.txt"},
			},
			expected: "File Tree:\n└── folder1\n    ├── file1.txt\n    └── file2.txt\n└── folder2\n    ├── file3.txt\n    └── file4.txt\n└── folder3\n    └── file5.txt\n",
			wantErr:  false,
		},
		{
			name: "Folders with paths",
			tree: map[string][]string{
				"folder1/path/to":    {"file1.txt"},
				"folder1/another/to": {"file2.txt"},
			},
			expected: "File Tree:\n└── folder1/another/to\n    └── file2.txt\n└── folder1/path/to\n    └── file1.txt\n",
			wantErr:  false,
		},
		{
			name: "Folders in unexpected order",
			tree: map[string][]string{
				"folder2": {"file3.txt"},
				"folder1": {"file1.txt", "file2.txt"},
				"folder3": {"file4.txt"},
			},
			expected: "File Tree:\n└── folder1\n    ├── file1.txt\n    └── file2.txt\n└── folder2\n    └── file3.txt\n└── folder3\n    └── file4.txt\n",
			wantErr:  false,
		},
		{
			name: "Nested folders representation",
			tree: map[string][]string{
				"parent":                  {"file1.txt"},
				"parent/child":            {"file2.txt", "file3.txt"},
				"parent/child/grandchild": {"file4.txt"},
			},
			expected: "File Tree:\n└── parent\n    └── file1.txt\n└── parent/child\n    ├── file2.txt\n    └── file3.txt\n└── parent/child/grandchild\n    └── file4.txt\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateFileTree(tt.tree)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFileTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("GenerateFileTree() = %v, want %v", got, tt.expected)

				// Print difference for debugging
				t.Errorf("Got:\n%s\nExpected:\n%s", strings.ReplaceAll(got, " ", "·"),
					strings.ReplaceAll(tt.expected, " ", "·"))
			}
		})
	}
}
