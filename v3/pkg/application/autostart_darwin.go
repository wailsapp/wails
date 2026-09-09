//go:build darwin && !ios && !server

package application

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/mac"
)

type darwinAutostart struct {
	app *App
}

func newAutostartImpl(app *App) autostartImpl {
	return &darwinAutostart{app: app}
}

// strategy picks SMAppService when running from a bundled .app on macOS 13+,
// otherwise the LaunchAgent plist path. Both paths support unbundled binaries
// (LaunchAgent path) so xbar-style scripts in development still work.
func (a *darwinAutostart) strategy() AutostartStrategy {
	if !runningFromAppBundle() {
		return AutostartStrategyLaunchAgent
	}
	if mac.GetBundleID() == "" {
		return AutostartStrategyLaunchAgent
	}
	major, _ := darwinMajorVersion()
	if major < 13 {
		return AutostartStrategyLaunchAgent
	}
	return AutostartStrategySMAppService
}

func (a *darwinAutostart) enable(opts AutostartOptions) error {
	if err := validateAutostartIdentifier(opts.Identifier); err != nil {
		return err
	}
	switch a.strategy() {
	case AutostartStrategySMAppService:
		if err := smAppServiceRegister(); err == nil {
			return nil
		} else if !errors.Is(err, errSMAppServiceUnavailable) {
			return fmt.Errorf("SMAppService register: %w", err)
		}
		fallthrough
	default:
		return a.enableLaunchAgent(opts)
	}
}

func (a *darwinAutostart) disable() error {
	// Try both paths and merge errors — a previous version may have used
	// the other strategy.
	var errs []error
	if a.strategy() == AutostartStrategySMAppService {
		if err := smAppServiceUnregister(); err != nil && !errors.Is(err, errSMAppServiceUnavailable) && !errors.Is(err, errSMAppServiceNotRegistered) {
			errs = append(errs, fmt.Errorf("SMAppService unregister: %w", err))
		}
	}
	if err := a.disableLaunchAgent(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (a *darwinAutostart) status() (AutostartStatus, error) {
	if a.strategy() == AutostartStrategySMAppService {
		enabled, err := smAppServiceIsEnabled()
		// errSMAppServiceRequiresApproval means the user disabled the
		// login item in System Settings — semantically that's "not
		// enabled", not a hard error. Treat it like Unavailable so the
		// LaunchAgent fallback still gets a chance to find a legacy
		// entry from before the app was bundled.
		if err != nil &&
			!errors.Is(err, errSMAppServiceUnavailable) &&
			!errors.Is(err, errSMAppServiceRequiresApproval) {
			return AutostartStatus{}, fmt.Errorf("SMAppService status: %w", err)
		}
		if enabled {
			return AutostartStatus{
				Enabled:  true,
				Path:     mac.GetBundleID(),
				Strategy: AutostartStrategySMAppService,
			}, nil
		}
	}
	// LaunchAgent path: also checked when SMAppService said no, so a
	// previously-registered LaunchAgent doesn't disappear from view after an
	// upgrade to a bundled build.
	path, ok, err := a.findLaunchAgent()
	if err != nil {
		return AutostartStatus{}, err
	}
	if ok {
		return AutostartStatus{
			Enabled:  true,
			Path:     path,
			Strategy: AutostartStrategyLaunchAgent,
		}, nil
	}
	return AutostartStatus{}, nil
}

func (a *darwinAutostart) launchAgentsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("autostart: %w", err)
	}
	return filepath.Join(home, "Library", "LaunchAgents"), nil
}

func (a *darwinAutostart) defaultLabel() string {
	if id := mac.GetBundleID(); id != "" {
		return id
	}
	return "wails.autostart." + autostartSlug(a.app.options.Name)
}

func (a *darwinAutostart) enableLaunchAgent(opts AutostartOptions) error {
	exe, err := resolvedExecutable()
	if err != nil {
		return err
	}
	dir, err := a.launchAgentsDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("autostart: create LaunchAgents dir: %w", err)
	}
	label := opts.Identifier
	if label == "" {
		label = a.defaultLabel()
	}
	path := filepath.Join(dir, label+".plist")
	body, err := launchAgentPlist(label, exe, opts.Arguments)
	if err != nil {
		return err
	}
	if err := writeFileAtomic(path, body, 0644); err != nil {
		return fmt.Errorf("autostart: write plist: %w", err)
	}
	// Best-effort: activate immediately for the current GUI session.
	_ = launchctlBootstrap(path)
	return nil
}

func (a *darwinAutostart) disableLaunchAgent() error {
	path, ok, err := a.findLaunchAgent()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	_ = launchctlBootout(path)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("autostart: remove plist: %w", err)
	}
	return nil
}

// findLaunchAgent looks for a plist in ~/Library/LaunchAgents whose
// ProgramArguments first element equals the current executable.
func (a *darwinAutostart) findLaunchAgent() (string, bool, error) {
	dir, err := a.launchAgentsDir()
	if err != nil {
		return "", false, err
	}
	exe, err := resolvedExecutable()
	if err != nil {
		return "", false, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("autostart: read LaunchAgents dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".plist") {
			continue
		}
		full := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		if plistFirstProgramArg(data) == exe {
			return full, true, nil
		}
	}
	return "", false, nil
}

// runningFromAppBundle reports whether the current executable lives inside a
// .app bundle (path ends with .app/Contents/MacOS/<name>).
func runningFromAppBundle() bool {
	exe, err := resolvedExecutable()
	if err != nil {
		return false
	}
	macOSDir := filepath.Dir(exe)
	contentsDir := filepath.Dir(macOSDir)
	appDir := filepath.Dir(contentsDir)
	return filepath.Base(macOSDir) == "MacOS" &&
		filepath.Base(contentsDir) == "Contents" &&
		strings.HasSuffix(appDir, ".app")
}

func darwinMajorVersion() (int, error) {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return 0, err
	}
	ver := strings.TrimSpace(string(out))
	if i := strings.IndexByte(ver, '.'); i > 0 {
		ver = ver[:i]
	}
	return strconv.Atoi(ver)
}

// launchctlBootstrap loads a plist into the current GUI session. Best effort —
// errors are ignored (the plist will still be picked up at next login).
//
// Indirected through a package-level variable so unit tests can replace it
// with a no-op: a test plist with RunAtLoad=true that successfully bootstraps
// would respawn the test binary recursively.
var launchctlBootstrap = func(plistPath string) error {
	target := fmt.Sprintf("gui/%d", os.Getuid())
	return exec.Command("launchctl", "bootstrap", target, plistPath).Run()
}

var launchctlBootout = func(plistPath string) error {
	target := fmt.Sprintf("gui/%d", os.Getuid())
	return exec.Command("launchctl", "bootout", target, plistPath).Run()
}

// --- plist marshalling ------------------------------------------------------

func launchAgentPlist(label, exe string, args []string) ([]byte, error) {
	progArgs := append([]string{exe}, args...)
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n")
	sb.WriteString(`<plist version="1.0">` + "\n")
	sb.WriteString("  <dict>\n")
	sb.WriteString("    <key>Label</key>\n")
	sb.WriteString("    <string>" + xmlEscape(label) + "</string>\n")
	sb.WriteString("    <key>ProgramArguments</key>\n")
	sb.WriteString("    <array>\n")
	for _, a := range progArgs {
		sb.WriteString("      <string>" + xmlEscape(a) + "</string>\n")
	}
	sb.WriteString("    </array>\n")
	sb.WriteString("    <key>RunAtLoad</key>\n")
	sb.WriteString("    <true/>\n")
	sb.WriteString("    <key>KeepAlive</key>\n")
	sb.WriteString("    <false/>\n")
	sb.WriteString("  </dict>\n")
	sb.WriteString("</plist>\n")
	return []byte(sb.String()), nil
}

func xmlEscape(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

// plistFirstProgramArg returns the first <string> element under the
// ProgramArguments array in a LaunchAgent plist. Empty string on any parse
// failure — we treat a malformed file as "not ours".
func plistFirstProgramArg(data []byte) string {
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	dec.Strict = false
	var inDict, inArray, captureKey bool
	var lastKey string
	for {
		tok, err := dec.Token()
		if err != nil {
			return ""
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "dict":
				inDict = true
			case "key":
				if inDict {
					captureKey = true
				}
			case "array":
				if lastKey == "ProgramArguments" {
					inArray = true
				}
			case "string":
				if inArray {
					var s string
					if err := dec.DecodeElement(&s, &t); err == nil {
						return s
					}
					return ""
				}
			}
		case xml.CharData:
			if captureKey {
				lastKey = string(t)
				captureKey = false
			}
		case xml.EndElement:
			if t.Name.Local == "array" && inArray {
				return ""
			}
		}
	}
}
