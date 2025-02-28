package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/MrZoidberg/contexify/pkg/os/mocks"
)

func TestCalculateFileHeaderSize(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected int64
	}{
		{
			name:     "Empty path",
			path:     "",
			expected: int64(len(getFileHeader(""))),
		},
		{
			name:     "Simple path",
			path:     "test.go",
			expected: int64(len(getFileHeader("test.go"))),
		},
		{
			name:     "Path with directory",
			path:     "dir/test.go",
			expected: int64(len(getFileHeader("dir/test.go"))),
		},
		{
			name:     "Complex path",
			path:     "some/deep/nested/directory/structure/file.go",
			expected: int64(len(getFileHeader("some/deep/nested/directory/structure/file.go"))),
		},
		{
			name:     "Path with backslashes",
			path:     "folder\\subfolder\\file.txt",
			expected: int64(len(getFileHeader("folder\\subfolder\\file.txt"))),
		},
		{
			name:     "Path with special characters",
			path:     "test-file_name.txt!",
			expected: int64(len(getFileHeader("test-file_name.txt!"))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateFileHeaderSize(tt.path)
			if got != tt.expected {
				t.Errorf("calculateFileHeaderSize(%q) = %d, want %d", tt.path, got, tt.expected)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		data      []byte
		delimiter string
		offset    int64
		wantN     int64
		wantErr   bool
		writeErrs []error
	}{
		{
			name:      "Simple write",
			path:      "test.txt",
			data:      []byte("test data"),
			delimiter: "\n",
			offset:    0,
			wantN:     int64(len(getFileHeader("test.txt")) + len("test data") + 1),
			wantErr:   false,
		},
		{
			name:      "Write with offset",
			path:      "test.txt",
			data:      []byte("test data"),
			delimiter: "\n",
			offset:    100,
			wantN:     int64(len(getFileHeader("test.txt")) + len("test data") + 1),
			wantErr:   false,
		},
		{
			name:      "Write without delimiter",
			path:      "test.txt",
			data:      []byte("test data"),
			delimiter: "",
			offset:    0,
			wantN:     int64(len(getFileHeader("test.txt")) + len("test data")),
			wantErr:   false,
		},
		{
			name:      "Error on header write",
			path:      "test.txt",
			data:      []byte("test data"),
			delimiter: "\n",
			offset:    0,
			writeErrs: []error{fmt.Errorf("header write error")},
			wantErr:   true,
		},
		{
			name:      "Error on data write",
			path:      "test.txt",
			data:      []byte("test data"),
			delimiter: "\n",
			offset:    0,
			writeErrs: []error{nil, fmt.Errorf("data write error")},
			wantErr:   true,
		},
		{
			name:      "Error on delimiter write",
			path:      "test.txt",
			data:      []byte("test data"),
			delimiter: "\n",
			offset:    0,
			writeErrs: []error{nil, nil, fmt.Errorf("delimiter write error")},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFile := &mocks.MockFile{}
			writeCallIndex := 0

			mockFile.WriteAtFunc = func(b []byte, off int64) (int, error) {
				var err error
				if tt.writeErrs != nil && writeCallIndex < len(tt.writeErrs) {
					err = tt.writeErrs[writeCallIndex]
				}
				writeCallIndex++
				return len(b), err
			}

			got, err := writeFile(tt.path, &tt.data, tt.delimiter, mockFile, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.wantN {
				t.Errorf("writeFile() = %v, want %v", got, tt.wantN)
			}
		})
	}
}

func TestProcess_Success(t *testing.T) {
	// Create a temporary directory for test files.
	tempDir, err := os.MkdirTemp("", "process_success")
	if err != nil {
		t.Fatalf("MkdirTemp error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create two temporary files with known content.
	file1Path := tempDir + "/file1.txt"
	content1 := []byte("hello world")
	if err := os.WriteFile(file1Path, content1, 0o644); err != nil {
		t.Fatalf("WriteFile error for file1: %v", err)
	}

	file2Path := tempDir + "/file2.txt"
	content2 := []byte("Go testing!")
	if err := os.WriteFile(file2Path, content2, 0o644); err != nil {
		t.Fatalf("WriteFile error for file2: %v", err)
	}

	// Create a temporary output file.
	outputPath := tempDir + "/output.txt"
	// Ensure output file does not exist yet.
	_ = os.Remove(outputPath)

	// Use a delimiter.
	delimiter := "\n"
	// Use writeTree as false to avoid relying on GenerateFileTree.
	paths := []string{file1Path, file2Path}

	// Call process.
	estimateTokensFunc := func(text, method string) (int, error) {
		return 10, nil
	}
	result, err := process(paths, false, outputPath, delimiter, estimateTokensFunc)
	if err != nil {
		t.Fatalf("process error: %v", err)
	}

	// Calculate expected totalSize.
	// For each file, expected bytes = file header size + file data length + delimiter length.
	var expectedSize int64
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil {
			t.Fatalf("Stat error for %s: %v", p, err)
		}
		expectedSize += fi.Size() + calculateFileHeaderSize(p) + int64(len(delimiter))
	}
	if result.totalSize != expectedSize {
		t.Errorf("process totalSize = %d, want %d", result.totalSize, expectedSize)
	}

	// Optionally, verify that output file contains expected content.
	outData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile error for output: %v", err)
	}
	// For each input file, output should contain its header, its data and the delimiter.
	for _, p := range paths {
		header := getFileHeader(p)
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("ReadFile error for %s: %v", p, err)
		}
		expectedFragment := header + string(data) + delimiter
		if !contains(string(outData), expectedFragment) {
			t.Errorf("output file missing expected fragment for %s", p)
		}
	}
}

func TestProcess_Error(t *testing.T) {
	// Use a path that does not exist.
	nonExistentPath := "nonexistent.txt"
	// Create a temporary file for output.
	tempOut, err := os.CreateTemp("", "process_error")
	if err != nil {
		t.Fatalf("CreateTemp error: %v", err)
	}
	outputPath := tempOut.Name()
	tempOut.Close()
	defer os.Remove(outputPath)

	// Call process with a non-existent file.
	estimateTokensFunc := func(text, method string) (int, error) {
		return 10, nil
	}
	_, err = process([]string{nonExistentPath}, false, outputPath, "\n", estimateTokensFunc)
	if err == nil {
		t.Errorf("process error = nil, want an error due to nonexistent input file")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (s != "" && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := range s {
		if len(s[i:]) >= len(substr) && s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
