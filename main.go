package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MrZoidberg/contexify/app"
	"github.com/MrZoidberg/contexify/pkg/log"
	flags "github.com/jessevdk/go-flags"
)

type Opts struct {
	Verbose           bool                 `short:"v" long:"verbose" description:"Show verbose debug information"`
	Input             string               `short:"i" long:"input" description:"Input folder path" required:"true" default:"."`
	Output            string               `short:"o" long:"output" description:"Output file path" required:"true" default:"context.txt"`
	Include           string               `long:"include" description:"File patterns to include separated by semicolon"`
	Exclude           string               `long:"exclude" description:"File patterns to exclude separated by semicolon" default:"LICENSE;CHANGELOG.md"`
	ConfigFile        func(s string) error `short:"c" long:"config" description:"Config file path inside Input folder" default:".contexify.yml"`
	DisableGitignore  bool                 `short:"g" long:"disable-gitignore" description:"Disable usage of .gitignore file to exclude files"`
	DisableFolderTree bool                 `long:"disable-folder-tree" description:"Do not add folder tree to the context"`
	NotRecursive      bool                 `long:"non-recursive" description:"Do not include subfolders"`
	Delimiter         string               `long:"delimiter" description:"Delimiter between files in output" default:"\n---\n"`

	Tokenizer struct {
		Skip bool `long:"skip" description:"skip calculating token count"`
	} `group:"tokenizer" namespace:"tokenizer"`
}

const hardIgnore = ".git/*;.gitignore;.vscode/*;.contexify.yml"

var opts Opts

var revision = "local"

func main() {
	// Parse command line arguments
	p := flags.NewParser(&opts, flags.Default)
	opts.ConfigFile = func(s string) error {
		configPath := filepath.Join(opts.Input, s)
		i := flags.NewIniParser(p)
		return i.ParseFile(configPath)
	}

	_, err := p.Parse()
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	log.SetupLog(opts.Verbose, false)

	log.Infof("Contexify %s\n", revision)
	log.Debugf("Options: %+v", opts)

	exclude := strings.Join([]string{hardIgnore, opts.Exclude}, ";")
	opts.Exclude = exclude

	err = app.Run(app.RunOptions{
		Input:            opts.Input,
		Output:           opts.Output,
		Include:          strings.Split(opts.Include, ";"),
		Exclude:          strings.Split(opts.Exclude, ";"),
		DisableGitignore: opts.DisableGitignore,
		HideTree:         opts.DisableFolderTree,
		NotRecursive:     opts.NotRecursive,
		Delimiter:        opts.Delimiter,
		Tokenizer: app.TokenizerOptions{
			Skip: opts.Tokenizer.Skip,
		},
	})

	if err != nil {
		os.Exit(1)
	}
}
