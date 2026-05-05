package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leaanthony/clir"

	"github.com/samber/lo"
)

func main() {
	app := clir.NewCli("sed", "A simple sed replacement", "v1")
	app.NewSubCommandFunction("replace", "Replace a string in files", ReplaceInFiles)
	err := app.Run()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

type ReplaceInFilesOptions struct {
	Dir        string `name:"dir" help:"Directory to search in"`
	OldString  string `name:"old" description:"The string to replace"`
	NewString  string `name:"new" description:"The string to replace with"`
	Extensions string `name:"ext" description:"The file extensions to process"`
	Ignore     string `name:"ignore" description:"The files to ignore"`
}

func ReplaceInFiles(options *ReplaceInFilesOptions) error {
	extensions := strings.Split(options.Extensions, ",")
	ignore := strings.Split(options.Ignore, ",")
	err := filepath.Walk(options.Dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if !lo.Contains(extensions, ext) {
			println("Skipping", path)
			return nil
		}
		filename := filepath.Base(path)
		if lo.Contains(ignore, filename) {
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
		return fmt.Errorf("Error while replacing in files: %v", err)
	}

	return nil
}
