//go:build ignore

// Generate release-local notes, or synchronize a tagged release into the archive.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v3/tasks/changelog"
)

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: generate-release-changelog generate|sync [flags]")
	}
	flags := flag.NewFlagSet(os.Args[1], flag.ContinueOnError)
	var options changelog.Options
	var tag string
	switch os.Args[1] {
	case "generate":
		flags.StringVar(&options.Repository, "repo", os.Getenv("GITHUB_REPOSITORY"), "owner/repository")
		flags.StringVar(&options.Target, "target", "HEAD", "tested commit")
		flags.StringVar(&options.Channel, "channel", "nightly", "nightly or weekly")
		flags.StringVar(&options.Version, "version", "", "release version")
		flags.StringVar(&options.Since, "since", "", "explicit ancestor baseline for bootstrapping")
	case "sync":
		flags.StringVar(&tag, "tag", "", "release tag")
	default:
		return fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err := flags.Parse(os.Args[2:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return err
	}
	root := strings.TrimSpace(string(output))
	if os.Args[1] == "sync" {
		return changelog.SyncArchive(root, tag)
	}
	receipt, err := changelog.Generate(root, options, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Generated %d entries from %s to %s\n", len(receipt.Entries), receipt.BaselineRef, receipt.Target)
	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
