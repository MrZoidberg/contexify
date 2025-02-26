package app

import (
	"path/filepath"

	realOS "os"

	"github.com/MrZoidberg/contexify/pkg/log"
	"github.com/denormal/go-gitignore"
)

func loadGitignore(folder string) (gitignore.GitIgnore, error) {
	ospath := filepath.Join(folder, ".gitignore")
	ignore, err := gitignore.NewFromFile(ospath)
	if err != nil {
		return nil, err
	}

	return ignore, nil
}

type traverser struct {
	folder       string
	useGitIgnore bool
	include      []string
	exclude      []string
	recursive    bool
}

func newTraverser(folder string, useGitIgnore, recursive bool, include, exclude []string) *traverser {
	return &traverser{
		folder:       folder,
		useGitIgnore: useGitIgnore,
		include:      include,
		exclude:      exclude,
		recursive:    recursive,
	}
}

type inputFile struct {
	path string
	ext  string
	size int64
}

func (t *traverser) Traverse() ([]inputFile, error) {
	var files []inputFile
	var gitIgnore gitignore.GitIgnore
	var err error

	if t.useGitIgnore {
		gitIgnore, err = loadGitignore(t.folder)
		if err != nil {
			return nil, err
		}
	}

	err = filepath.Walk(t.folder, func(path string, info realOS.FileInfo, err error) error {
		if err != nil {
			log.Errorf("error accessing path %q: %v", path, err)
			return err
		}

		path = filepath.ToSlash(path)

		// Skip directories unless we're processing them
		if info.IsDir() {
			if !t.recursive && path != t.folder {
				log.Debugf("skipping directory %q", path)
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file should be ignored by .gitignore
		if t.useGitIgnore {
			if gitIgnore.Ignore(path) {
				log.Debugf("skipping file %q by .gitignore", path)
				return nil
			}
		}

		// Check include patterns using filepath.Match()
		if len(t.include) > 0 {
			matched := false
			for _, pattern := range t.include {
				if pattern == "" {
					matched = true
					break
				}
				if match, _ := filepath.Match(pattern, path); match {
					matched = true
					break
				}
			}
			if !matched {
				log.Debugf("skipping file %q by include patterns", path)
				return nil
			}
		}

		// Check exclude patterns using filepath.Match()
		for _, pattern := range t.exclude {
			if pattern == "" {
				continue
			}
			if match, _ := filepath.Match(pattern, path); match {
				log.Debugf("skipping file %q by exclude pattern %q", path, pattern)
				return nil
			}
		}

		log.Debugf("adding file %q", path)
		files = append(files, inputFile{
			path: path,
			ext:  filepath.Ext(path),
			size: info.Size(),
		})
		return nil
	})

	return files, err
}
