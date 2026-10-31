// Package installerbuild prepares and publishes NSIS installers while retaining
// the project's installer script and the NSIS compiler as a tool boundary.
package installerbuild

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/buildmeta"
)

type File struct {
	Name string
	Path string
}

type Options struct {
	Script       string
	Output       string
	Arch         string
	Scope        string
	Primary      File
	Companions   []File
	Version      *buildmeta.Info
	Prepare      func(context.Context, string) error
	Defines      map[string]string
	SkipWebView2 bool
}

type Compile func(ctx context.Context, directory string, args []string) error

func (o Options) Validate() error {
	if err := validateDefines(o.Defines); err != nil {
		return err
	}
	if o.Arch != "amd64" && o.Arch != "arm64" {
		return fmt.Errorf("NSIS requires amd64 or arm64, got %q", o.Arch)
	}
	if o.Scope != "" && o.Scope != "user" && o.Scope != "machine" {
		return fmt.Errorf("NSIS install scope must be user or machine, got %q", o.Scope)
	}
	if o.Script == "" || o.Output == "" {
		return fmt.Errorf("NSIS script and output are required")
	}
	for _, path := range []string{o.Script, o.Output} {
		if err := validateSourcePath(path); err != nil {
			return err
		}
	}
	if err := validateFiles(append([]File{o.Primary}, o.Companions...)); err != nil {
		return err
	}
	if o.Version != nil {
		_, err := o.Version.Windows()
		return err
	}
	return nil
}

func validateFiles(files []File) error {
	seen := make(map[string]bool, len(files))
	for _, file := range files {
		if file.Name == "" || file.Name == "." || strings.HasSuffix(file.Name, ".") || strings.HasSuffix(file.Name, " ") || strings.ContainsAny(file.Name, "<>:/\\|?*\x00\r\n\"") || strings.EqualFold(file.Name, "uninstall.exe") {
			return fmt.Errorf("invalid installer executable name %q", file.Name)
		}
		if file.Path == "" || seen[strings.ToLower(file.Name)] {
			return fmt.Errorf("missing or duplicate installer executable %q", file.Name)
		}
		if err := validateSourcePath(file.Path); err != nil {
			return err
		}
		seen[strings.ToLower(file.Name)] = true
	}
	return nil
}

// Create overrides only the generated installer's destination. Relative script
// assets continue to resolve beside the project's script. A private wrapper
// adds companion executables, leaving the project's script and binary untouched.
func Create(ctx context.Context, options Options, compile Compile) error {
	if err := options.Validate(); err != nil {
		return err
	}
	if compile == nil {
		return fmt.Errorf("NSIS compiler is unavailable")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	script, err := filepath.Abs(options.Script)
	if err != nil {
		return err
	}
	output, err := filepath.Abs(options.Output)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(filepath.Dir(output), ".wake-installer-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	stagedOutput := filepath.Join(stage, "installer.exe")
	wrapper, err := wrapperScript(script, stagedOutput, options)
	if err != nil {
		return err
	}
	path := filepath.Join(stage, "wake.nsi")
	if err = os.WriteFile(path, []byte(wrapper), 0600); err != nil {
		return err
	}
	if err = compile(ctx, filepath.Dir(script), []string{"-NOCD", "-INPUTCHARSET", "UTF8", path}); err != nil {
		return fmt.Errorf("compile NSIS installer: %w", err)
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	if err = options.prepare(ctx, stagedOutput); err != nil {
		return err
	}
	return os.Rename(stagedOutput, output)
}

func (options Options) prepare(ctx context.Context, path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("NSIS compiler output: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() == 0 {
		return fmt.Errorf("NSIS compiler did not produce a regular, non-empty installer")
	}
	if options.Prepare != nil {
		if err = options.Prepare(ctx, path); err != nil {
			return err
		}
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return nil
}

func wrapperScript(script, output string, options Options) (string, error) {
	var text strings.Builder
	define := func(name, value string) { fmt.Fprintf(&text, "!define %s %s\n", name, quote(value)) }
	keys := make([]string, 0, len(options.Defines))
	for name := range options.Defines {
		keys = append(keys, name)
	}
	slices.Sort(keys)
	for _, name := range keys {
		define(name, options.Defines[name])
	}
	primary, err := filepath.Abs(options.Primary.Path)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(&text, "!define ARG_WAILS_%s_BINARY %s\n", strings.ToUpper(options.Arch), sourceQuote(primary))
	define("PRODUCT_EXECUTABLE", options.Primary.Name)
	if options.SkipWebView2 {
		define("WAILS_SKIP_WEBVIEW2", "1")
	}
	if options.Scope != "" {
		define("WAILS_INSTALL_SCOPE", options.Scope)
		if options.Scope == "user" {
			define("REQUEST_EXECUTION_LEVEL", "user")
		}
	}
	if options.Version != nil {
		define("INFO_PRODUCTVERSION", options.Version.Core)
	}
	fmt.Fprintf(&text, "!include %s\n", sourceQuote(filepath.Base(script)))
	if len(options.Companions) > 0 {
		text.WriteString("Section \"-Wake companion executables\"\nSetOutPath \"$INSTDIR\"\n")
		for _, file := range options.Companions {
			path, err := filepath.Abs(file.Path)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&text, "File %s %s\n", quote("/oname="+file.Name), sourceQuote(path))
		}
		text.WriteString("SectionEnd\n")
	}
	fmt.Fprintf(&text, "OutFile %s\n", sourceQuote(output))
	return text.String(), nil
}

var runtimeEscapes = strings.NewReplacer("$", "$$", `"`, `$\"`, "\r", `$\r`, "\n", `$\n`)

func quote(value string) string { return `"` + runtimeEscapes.Replace(value) + `"` }

// Compile-time filesystem arguments do not decode NSIS runtime dollar escapes.
func sourceQuote(value string) string { return `"` + value + `"` }

func validateSourcePath(path string) error {
	if strings.ContainsAny(path, "\x00\r\n\"") || strings.Contains(path, "${") {
		return fmt.Errorf("installer source path cannot contain quotes, control characters or NSIS substitutions: %q", path)
	}
	return nil
}

func validateDefines(defines map[string]string) error {
	for name, value := range defines {
		if name == "" || strings.ContainsRune(value, 0) {
			return fmt.Errorf("invalid NSIS define %q", name)
		}
		for i, c := range name {
			if c != '_' && !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(i > 0 && c >= '0' && c <= '9') {
				return fmt.Errorf("invalid NSIS define %q", name)
			}
		}
		switch name {
		case "PRODUCT_EXECUTABLE", "WAILS_INSTALL_SCOPE", "REQUEST_EXECUTION_LEVEL", "INFO_PRODUCTVERSION", "WAILS_SKIP_WEBVIEW2":
			return fmt.Errorf("NSIS define %s is owned by Wake", name)
		}
		if strings.HasPrefix(name, "ARG_WAILS_") {
			return fmt.Errorf("NSIS define %s is owned by Wake", name)
		}
	}
	return nil
}
