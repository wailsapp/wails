package application

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/internal/fileexplorer"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

// EnvironmentManager manages environment-related operations
type EnvironmentManager struct {
	app *App
}

// newEnvironmentManager creates a new EnvironmentManager instance
func newEnvironmentManager(app *App) *EnvironmentManager {
	return &EnvironmentManager{
		app: app,
	}
}

// Info returns environment information
func (em *EnvironmentManager) Info() EnvironmentInfo {
	info, _ := operatingsystem.Info()
	result := EnvironmentInfo{
		OS:     runtime.GOOS,
		Arch:   runtime.GOARCH,
		Debug:  em.app.isDebugMode,
		OSInfo: info,
	}
	result.PlatformInfo = em.app.platformEnvironment()
	return result
}

// IsDarkMode returns true if the system is in dark mode
func (em *EnvironmentManager) IsDarkMode() bool {
	if em.app.impl == nil {
		return false
	}
	return em.app.impl.isDarkMode()
}

// GetAccentColor returns the system accent color
func (em *EnvironmentManager) GetAccentColor() string {
	if em.app.impl == nil {
		return "rgb(0,122,255)"
	}
	return em.app.impl.getAccentColor()
}

// OpenFileManager opens the file manager at the specified path, optionally selecting the file
func (em *EnvironmentManager) OpenFileManager(path string, selectFile bool) error {
	return InvokeSyncWithError(func() error {
		return fileexplorer.OpenFileManager(path, selectFile)
	})
}

func (em *EnvironmentManager) HasFocusFollowsMouse() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	info := em.app.platformEnvironment()
	if ffm, ok := info["focusFollowsMouse"].(bool); ok {
		return ffm
	}
	return false
}

// AccessibilitySettings mirrors the user's system accessibility preferences
// that affect how an app should draw and animate.
type AccessibilitySettings struct {
	// ReduceMotion asks apps to avoid non-essential animation.
	ReduceMotion bool `json:"reduceMotion"`
	// ReduceTransparency asks apps to use opaque backgrounds instead of
	// blur and vibrancy.
	ReduceTransparency bool `json:"reduceTransparency"`
	// IncreaseContrast asks for stronger borders and higher-contrast colours.
	IncreaseContrast bool `json:"increaseContrast"`
	// DifferentiateWithoutColor asks apps not to rely on colour alone to
	// convey meaning.
	DifferentiateWithoutColor bool `json:"differentiateWithoutColor"`
	// InvertColors reports that display colours are inverted system-wide.
	InvertColors bool `json:"invertColors"`
	// VoiceOverEnabled reports whether the VoiceOver screen reader is on.
	VoiceOverEnabled bool `json:"voiceOverEnabled"`
	// SwitchControlEnabled reports whether Switch Control is on.
	SwitchControlEnabled bool `json:"switchControlEnabled"`
}

// KeyboardLayout describes the active keyboard input source.
type KeyboardLayout struct {
	// ID is the platform identifier, for example
	// "com.apple.keylayout.British" on macOS.
	ID string `json:"id"`
	// Name is the user-visible, localised name ("British").
	Name string `json:"name"`
	// Languages lists the BCP 47 language tags the layout is intended for.
	Languages []string `json:"languages"`
}

// LocaleInfo describes the locale the application is running under.
type LocaleInfo struct {
	// Identifier is the full locale identifier, for example "en_GB".
	Identifier string `json:"identifier"`
	// Language is the ISO 639 language code ("en").
	Language string `json:"language"`
	// Region is the ISO 3166 region code ("GB"); empty when unknown.
	Region string `json:"region"`
	// Preferred is the user's ordered list of preferred languages as BCP 47
	// tags ("en-GB", "fr-FR"), independent of what the app supports.
	Preferred []string `json:"preferred"`
}

// Accessibility returns the user's accessibility display preferences.
// Changes are announced through events.Common.AccessibilitySettingsChanged
// (and events.Mac.ApplicationDidChangeAccessibilitySettings on macOS).
//
// macOS: NSWorkspace accessibilityDisplayShould* properties plus
// voiceOverEnabled and switchControlEnabled.
//
// Other platforms: all false.
func (em *EnvironmentManager) Accessibility() AccessibilitySettings {
	return platformAccessibilitySettings()
}

// KeyboardLayout returns the active keyboard input source. Changes are
// announced through events.Mac.ApplicationDidChangeKeyboardLayout.
//
// macOS: TISCopyCurrentKeyboardInputSource.
//
// Other platforms: zero value.
func (em *EnvironmentManager) KeyboardLayout() KeyboardLayout {
	return platformKeyboardLayout()
}

// Locale returns the locale the application runs under. Changes are
// announced through events.Mac.ApplicationDidChangeLocale.
//
// macOS: NSLocale currentLocale and preferredLanguages. Note that
// currentLocale is the locale AppKit selected for the app, not simply the
// user's first preferred language: it is the first entry of the user's
// language list that the bundle declares support for through the
// CFBundleLocalizations Info.plist key (or .lproj directories), combined
// with the user's region. A bundle that declares no localisations is
// treated as English-only, so a German user's app reports "de_DE" only if
// Info.plist lists "de" in CFBundleLocalizations (see issue #4582).
// Preferred is unaffected and always carries the user's full ordered list,
// so use it to pick a language yourself when the bundle cannot declare
// every localisation up front.
//
// Other platforms: derived from LC_ALL, LC_MESSAGES and LANG when set,
// otherwise the zero value.
func (em *EnvironmentManager) Locale() LocaleInfo {
	return platformLocale()
}

// parsePosixLocale converts a POSIX locale string such as "en_GB.UTF-8" or
// "pt_BR" into LocaleInfo. It reports false for the "C" and "POSIX" locales
// and for empty input.
func parsePosixLocale(value string) (LocaleInfo, bool) {
	value = strings.TrimSpace(value)
	if dot := strings.IndexAny(value, ".@"); dot >= 0 {
		value = value[:dot]
	}
	if value == "" || value == "C" || value == "POSIX" {
		return LocaleInfo{}, false
	}
	info := LocaleInfo{Identifier: value}
	language, region, _ := strings.Cut(value, "_")
	info.Language = language
	info.Region = region
	tag := language
	if region != "" {
		tag += "-" + region
	}
	info.Preferred = []string{tag}
	return info, true
}

// localeFromEnvironment applies parsePosixLocale to the usual environment
// variables in precedence order.
func localeFromEnvironment(lookup func(string) string) LocaleInfo {
	for _, name := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if info, ok := parsePosixLocale(lookup(name)); ok {
			return info
		}
	}
	return LocaleInfo{}
}

// DefaultHandler selects what a default-application query or change refers
// to: exactly one of ContentType or URLScheme must be set.
type DefaultHandler struct {
	// ContentType is a uniform type identifier such as "public.plain-text"
	// or "com.adobe.pdf".
	ContentType string `json:"contentType,omitempty"`
	// URLScheme is a URL scheme without the colon, such as "mailto" or
	// "myapp".
	URLScheme string `json:"urlScheme,omitempty"`
}

func (h DefaultHandler) validate() error {
	contentType := strings.TrimSpace(h.ContentType)
	scheme := strings.TrimSpace(h.URLScheme)
	switch {
	case contentType == "" && scheme == "":
		return errors.New("environment: DefaultHandler needs a ContentType or a URLScheme")
	case contentType != "" && scheme != "":
		return errors.New("environment: DefaultHandler takes either a ContentType or a URLScheme, not both")
	case contentType != "" && strings.ContainsAny(contentType, " \t\r\n/:"):
		return fmt.Errorf("environment: %q is not a uniform type identifier", h.ContentType)
	case scheme != "":
		for index, r := range scheme {
			isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
			isDigit := r >= '0' && r <= '9'
			if !(isLetter || (index > 0 && (isDigit || r == '+' || r == '-' || r == '.'))) {
				return fmt.Errorf("environment: %q is not a URL scheme (letters, digits, '+', '-' and '.', starting with a letter, no colon)", h.URLScheme)
			}
		}
	}
	return nil
}

// String describes the handler for error messages.
func (h DefaultHandler) String() string {
	if h.ContentType != "" {
		return "content type " + h.ContentType
	}
	return "URL scheme " + h.URLScheme
}

// ErrDefaultHandlerUnsupported is returned by SetDefaultHandler and
// DefaultHandler on platforms without a default-application registry.
var ErrDefaultHandlerUnsupported = errors.New("environment: default handlers are not supported on this platform")

// SetDefaultHandler makes this application the default for a content type
// or URL scheme. The application must be a bundle with a bundle identifier;
// unbundled `go run` binaries cannot be registered.
//
// macOS: NSWorkspace setDefaultApplicationAtURL:toOpenContentType: and
// toOpenURLsWithScheme: (macOS 12+), which may show a confirmation dialog
// for browsers and mail clients, with LSSetDefaultRoleHandlerForContentType
// and LSSetDefaultHandlerForURLScheme on older systems.
//
// Other platforms: ErrDefaultHandlerUnsupported.
func (em *EnvironmentManager) SetDefaultHandler(handler DefaultHandler) error {
	if err := handler.validate(); err != nil {
		return err
	}
	return platformSetDefaultHandler(handler)
}

// DefaultHandler returns the application that currently opens the content
// type or URL scheme.
//
// macOS: NSWorkspace URLForApplicationToOpenContentType: and
// URLForApplicationToOpenURL: (macOS 12+), with the Launch Services
// LSCopyDefaultRoleHandlerForContentType and LSCopyDefaultHandlerForURLScheme
// on older systems. ErrApplicationNotFound when nothing is registered.
//
// Other platforms: ErrDefaultHandlerUnsupported.
func (em *EnvironmentManager) DefaultHandler(handler DefaultHandler) (AppInfo, error) {
	if err := handler.validate(); err != nil {
		return AppInfo{}, err
	}
	return platformDefaultHandler(handler)
}
