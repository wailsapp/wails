package buildwarnings

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/internal/term"
)

// EnvVar is the environment variable that holds the path to the warnings file.
// It is set by the build command and inherited by all subprocess commands so
// they can append warnings without knowing the parent PID or session ID.
const EnvVar = "WAILS_WARNINGS_FILE"

// Add appends a warning to the current build session's warnings file.
// source identifies where the warning comes from (e.g. "tool has-cc").
// If EnvVar is not set (i.e. we are not inside a build), this is a no-op.
func Add(source, message string) {
	path := os.Getenv(EnvVar)
	if path == "" {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s\t%s\n", source, message)
}

// FlushAndPrint reads all warnings from the current build session's warnings
// file, prints them using the terminal warning printer, then removes the file.
// Safe to call when no file exists or EnvVar is unset.
func FlushAndPrint() {
	path := os.Getenv(EnvVar)
	if path == "" {
		return
	}
	warnings := read(path)
	_ = os.Remove(path)
	if len(warnings) == 0 {
		return
	}
	term.Warning("Build warnings:")
	for _, w := range warnings {
		term.Warningf("  [%s] %s", w.source, w.message)
	}
}

type entry struct {
	source  string
	message string
}

func read(path string) []entry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []entry
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		source, message, ok := strings.Cut(line, "\t")
		if !ok {
			continue
		}
		out = append(out, entry{source: source, message: message})
	}
	return out
}
