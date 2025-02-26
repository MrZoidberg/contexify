package app

import (
	"fmt"
	realOS "os"
	"path/filepath"
	"sync"
	"time"

	"github.com/MrZoidberg/contexify/pkg/os"

	"github.com/MrZoidberg/contexify/pkg/log"
	"github.com/dustin/go-humanize"
)

// TokenizerOptions contains options for the tokenizer
type TokenizerOptions struct {
	// Skip skips calculating token count
	Skip bool
}

// RunOptions contains options for the Run function
type RunOptions struct {
	// Input folder path
	Input string
	// Output file path
	Output string
	// Include file patterns separated by semicolon
	Include []string
	// Exclude file patterns separated by semicolon
	Exclude []string
	// DisableGitignore disables usage of .gitignore file to exclude files
	DisableGitignore bool
	// HideTree disables adding folder tree to the context
	HideTree bool
	// NotRecursive disables including subfolders
	NotRecursive bool
	// Delimiter between files in output
	Delimiter string
	// Tokenizer options
	Tokenizer TokenizerOptions
}

func getFileHeader(path string) string {
	return fmt.Sprintf("===> File: %s\n\n", filepath.ToSlash(path))
}

func calculateFileHeaderSize(path string) int64 {
	getFileHeaderSize := int64(len(getFileHeader(path)))
	return getFileHeaderSize
}

func writeFile(path string, data *[]byte, delimiter string, file *realOS.File, offset int64) (int64, error) {
	header := getFileHeader(path)
	n1, err := file.WriteAt([]byte(header), offset)
	if err != nil {
		return 0, err
	}
	n2, err := file.WriteAt(*data, offset+int64(n1))
	if err != nil {
		return 0, err
	}
	if delimiter != "" {
		_, err = file.WriteAt([]byte(delimiter), offset+int64(n1)+int64(n2))
		if err != nil {
			return 0, err
		}
	}
	return int64(n1 + n2 + len(delimiter)), nil
}

type processingResult struct {
	totalSize   int64
	totalTokens int
}

func process(paths []string, writeTree bool, output, delimiter string) (processingResult, error) {
	currentOffset := int64(0)
	totalTokens := 0

	out, err := os.Create(output)
	if err != nil {
		return processingResult{}, err
	}
	defer out.Close()

	// map folder to files
	folderMap := make(map[string][]string)
	// map folder to offset in output file
	folderOffsets := make(map[string]int64)

	for _, p := range paths {
		folder := filepath.Dir(p)
		folderMap[folder] = append(folderMap[folder], p)
	}

	// write folder tree
	if writeTree {
		tree, err := GenerateFileTree(folderMap)
		if err != nil {
			return processingResult{}, fmt.Errorf("error generating file tree: %w", err)
		}
		_, err = out.WriteString(tree + delimiter)
		if err != nil {
			return processingResult{}, fmt.Errorf("error writing file tree: %w", err)
		}
		currentOffset += int64(len(tree) + len(delimiter))
	}

	for folder, files := range folderMap {
		folderOffsets[folder] = currentOffset
		for _, file := range files {
			fileInfo, err := os.Stat(file)
			if err != nil {
				return processingResult{}, err
			}
			currentOffset += fileInfo.Size() + calculateFileHeaderSize(file) + int64(len(delimiter))
		}
	}

	// write output
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 1)
	)

	for folder, files := range folderMap {
		wg.Add(1)
		go func(folder string, files []string) {
			defer wg.Done()
			// Get the starting offset for this folder.
			offset := folderOffsets[folder]
			for _, file := range files {
				data, err := os.ReadFile(file)
				if err != nil {
					errChan <- err
					return
				}
				// estimate tokens
				tokens, err := EstimateTokens(string(data), "max")
				if err != nil {
					errChan <- err
					return
				}
				totalTokens += tokens

				// write header and data at the given offset.
				n, err := writeFile(file, &data, delimiter, out, offset)
				if err != nil {
					errChan <- err
					return
				}
				// Move offset by the number of bytes written.
				offset += n
			}
		}(folder, files)
	}

	// Wait for all folder processing to complete.
	wg.Wait()

	// Check if any goroutine reported an error.
	select {
	case err := <-errChan:
		return processingResult{}, err
	default:
		// All writes succeeded.
	}

	return processingResult{
		totalSize:   currentOffset,
		totalTokens: totalTokens,
	}, nil
}

// Run runs the processing of the input folder and writes the output to the output file
func Run(options RunOptions) error {
	// Validate options
	if err := validateOptions(options); err != nil {
		return err
	}

	// Traverse input folder
	traverser := newTraverser(options.Input, !options.DisableGitignore, !options.NotRecursive, options.Include, options.Exclude)
	files, err := traverser.Traverse()
	if err != nil {
		log.Errorf("error traversing folder %s: %v", options.Input, err)
		return fmt.Errorf("error traversing folder %s: %v", options.Input, err)
	}

	// Write stats on found files
	totalSize := int64(0)
	for _, file := range files {
		totalSize += file.size
	}
	log.Infof("Found %d files, total size: %s", len(files), humanize.Bytes(uint64(totalSize)))

	// Process files
	filePaths := make([]string, len(files))
	for i, file := range files {
		filePaths[i] = file.path
	}
	startTime := time.Now()

	result, err := process(filePaths, !options.HideTree, options.Output, options.Delimiter)
	if err != nil {
		log.Errorf("error processing files: %v", err)
		return fmt.Errorf("error processing files: %v", err)
	}

	processingTime := time.Since(startTime)
	log.Infof("Processed %d files, total size: %s, total tokens estimate: %d, processing time: %s",
		len(files),
		humanize.Bytes(uint64(result.totalSize)),
		result.totalTokens,
		processingTime)

	return nil
}

func validateOptions(options RunOptions) error {
	// Validate input folder
	if _, err := os.Stat(options.Input); os.IsNotExist(err) {
		log.Errorf("input folder %s does not exist", options.Input)
		return fmt.Errorf("input folder %s does not exist", options.Input)
	}

	return nil
}
