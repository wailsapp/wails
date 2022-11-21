package flags

import (
	"fmt"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/system"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
	"github.com/wailsapp/wails/v2/pkg/commands/buildtags"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	Quiet   int = 0
	Normal  int = 1
	Verbose int = 2
)

// TODO: unify this and `build.Options`
type Build struct {
	NoPackage               bool   `name:"noPackage" description:"Skips platform specific packaging"`
	Compiler                string `description:"Use a different go compiler to build, eg go1.15beta1"`
	SkipModTidy             bool   `name:"m" description:"Skip mod tidy before compile"`
	Upx                     bool   `description:"Compress final binary with UPX (if installed)"`
	UpxFlags                string `description:"Flags to pass to upx"`
	Platform                string `description:"Platform to target. Comma separate multiple platforms"`
	Verbosity               int    `name:"v" description:"Verbosity level (0 = quiet, 1 = normal, 2 = verbose)"`
	LdFlags                 string `description:"Additional ldflags to pass to the compiler"`
	Tags                    string `description:"Build tags to pass to Go compiler. Must be quoted. Space or comma (but not both) separated"`
	OutputFilename          string `name:"o" description:"Output filename"`
	Clean                   bool   `description:"Clean the bin directory before building"`
	WebView2                string `description:"WebView2 installer strategy: download,embed,browser,error"`
	SkipFrontend            bool   `name:"s" description:"Skips building the frontend"`
	ForceBuild              bool   `name:"f" description:"Force build of application"`
	UpdateWailsVersionGoMod bool   `name:"u" description:"Updates go.mod to use the same Wails version as the CLI"`
	Debug                   bool   `description:"Builds the application in debug mode"`
	NSIS                    bool   `description:"Generate NSIS installer for Windows"`
	TrimPath                bool   `description:"Remove all file system paths from the resulting executable"`
	RaceDetector            bool   `description:"Build with Go's race detector"`
	WindowsConsole          bool   `description:"Keep the console when building for Windows"`
	Obfuscated              bool   `description:"Code obfuscation of bound Wails methods"`
	GarbleArgs              string `description:"Arguments to pass to garble"`
	DryRun                  bool   `description:"Prints the build command without executing it"`
	SkipBindings            bool   `description:"Skips generation of bindings"`

	// Internal state
	compilerPath  string
	userTags      []string
	wv2rtstrategy string // WebView2 runtime strategy
	defaultArch   string // Default architecture
}

func (b *Build) Default() *Build {

	defaultPlatform := os.Getenv("GOOS")
	if defaultPlatform == "" {
		defaultPlatform = runtime.GOOS
	}
	defaultArch := os.Getenv("GOARCH")
	if defaultArch == "" {
		if system.IsAppleSilicon {
			defaultArch = "arm64"
		} else {
			defaultArch = runtime.GOARCH
		}
	}
	b.defaultArch = defaultArch
	platform := defaultPlatform + "/" + defaultArch

	return &Build{
		Compiler:   "go",
		Platform:   platform,
		Verbosity:  Normal,
		WebView2:   "download",
		GarbleArgs: "-literals -tiny -seed=random",
	}
}

func (b *Build) GetBuildMode() build.Mode {
	if b.Debug {
		return build.Debug
	}
	return build.Production
}

func (b *Build) GetWebView2Strategy() string {
	return b.wv2rtstrategy
}

func (b *Build) GetTargets() *slicer.StringSlicer {
	var targets slicer.StringSlicer
	targets.AddSlice(strings.Split(b.Platform, ","))
	targets.Deduplicate()
	return &targets
}

func (b *Build) GetCompilerPath() string {
	return b.compilerPath
}

func (b *Build) GetTags() []string {
	return b.userTags
}

func (b *Build) Process() error {
	// Lookup compiler path
	var err error
	b.compilerPath, err = exec.LookPath(b.Compiler)
	if err != nil {
		return fmt.Errorf("unable to find compiler: %s", b.Compiler)
	}

	// Process User Tags
	b.userTags, err = buildtags.Parse(b.Tags)
	if err != nil {
		return err
	}

	// WebView2 installer strategy (download by default)
	b.WebView2 = strings.ToLower(b.WebView2)
	if b.WebView2 != "" {
		validWV2Runtime := slicer.String([]string{"download", "embed", "browser", "error"})
		if !validWV2Runtime.Contains(b.WebView2) {
			return fmt.Errorf("invalid option for flag 'webview2': %s", b.WebView2)
		}
		b.wv2rtstrategy = "wv2runtime." + b.WebView2
	}

	return nil
}

func bool2Str(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func (b *Build) GetBuildModeAsString() string {
	if b.Debug {
		return "debug"
	}
	return "production"
}

func (b *Build) GetDefaultArch() string {
	return b.defaultArch
}

/*
	_, _ = fmt.Fprintf(w, "Frontend Directory: \t%s\n", projectOptions.GetFrontendDir())
	_, _ = fmt.Fprintf(w, "Obfuscated: \t%t\n", buildOptions.Obfuscated)
	if buildOptions.Obfuscated {
		_, _ = fmt.Fprintf(w, "Garble Args: \t%s\n", buildOptions.GarbleArgs)
	}
	_, _ = fmt.Fprintf(w, "Skip Frontend: \t%t\n", skipFrontend)
	_, _ = fmt.Fprintf(w, "Compress: \t%t\n", buildOptions.Compress)
	_, _ = fmt.Fprintf(w, "Package: \t%t\n", buildOptions.Pack)
	_, _ = fmt.Fprintf(w, "Clean Bin Dir: \t%t\n", buildOptions.CleanBinDirectory)
	_, _ = fmt.Fprintf(w, "LDFlags: \t\"%s\"\n", buildOptions.LDFlags)
	_, _ = fmt.Fprintf(w, "Tags: \t[%s]\n", strings.Join(buildOptions.UserTags, ","))
	_, _ = fmt.Fprintf(w, "Race Detector: \t%t\n", buildOptions.RaceDetector)
	if len(buildOptions.OutputFile) > 0 && targets.Length() == 1 {
		_, _ = fmt.Fprintf(w, "Output File: \t%s\n", buildOptions.OutputFile)
	}
*/
