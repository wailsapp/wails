package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func run(args []string, output io.Writer) error {
	root := flag.NewFlagSet("sed", flag.ContinueOnError)
	root.SetOutput(io.Discard)
	rootHelp := root.Bool("help", false, "Get help on the 'sed' command")
	err := root.Parse(args)
	if err != nil && !errors.Is(err, flag.ErrHelp) {
		return err
	}
	args = root.Args()
	if len(args) == 0 || *rootHelp || errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(output, "sed v1 - A simple sed replacement\n\nUsage: sed replace [flags]\nUse 'sed replace --help' for replacement options.")
		return nil
	}
	if args[0] != "replace" {
		return fmt.Errorf("unknown command %q; use 'sed --help' for usage", args[0])
	}
	var options ReplaceInFilesOptions
	flags := flag.NewFlagSet("sed replace", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	help := flags.Bool("help", false, "Get help on the 'sed replace' command")
	flags.StringVar(&options.Dir, "dir", "", "Directory to search in")
	flags.StringVar(&options.OldString, "old", "", "The string to replace")
	flags.StringVar(&options.NewString, "new", "", "The string to replace with")
	flags.StringVar(&options.Extensions, "ext", "", "The file extensions to process")
	flags.StringVar(&options.Ignore, "ignore", "", "The files to ignore")
	err = flags.Parse(args[1:])
	if err != nil && !errors.Is(err, flag.ErrHelp) {
		return err
	}
	if *help || errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(output, "sed v1 - A simple sed replacement\n\nsed replace - Replace a string in files\nFlags:")
		flags.SetOutput(output)
		flags.PrintDefaults()
		return nil
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}
	return ReplaceInFiles(&options)
}

type ReplaceInFilesOptions struct {
	Dir        string
	OldString  string
	NewString  string
	Extensions string
	Ignore     string
}

func ReplaceInFiles(options *ReplaceInFilesOptions) error {
	extensions := strings.Split(options.Extensions, ",")
	ignore := strings.Split(options.Ignore, ",")
	err := filepath.Walk(options.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if !slices.Contains(extensions, ext) {
			println("Skipping", path)
			return nil
		}
		filename := filepath.Base(path)
		if slices.Contains(ignore, filename) {
			println("Ignoring:", path)
			return nil
		}

		println("Processing file:", path)

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newContent := strings.Replace(string(content), options.OldString, options.NewString, -1)

		return os.WriteFile(path, []byte(newContent), info.Mode())
	})

	if err != nil {
		return fmt.Errorf("Error while replacing in files: %w", err)
	}

	return nil
}
