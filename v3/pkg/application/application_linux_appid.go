//go:build linux && cgo && !android && !server

package application

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var invalidAppNameChars = regexp.MustCompile(`[^a-zA-Z0-9_-]`)
var leadingDigits = regexp.MustCompile(`^[0-9]+`)

// sanitizeAppName sanitizes the application name into a single element of a
// GTK/D-Bus application id: only alphanumeric characters, hyphens and
// underscores, and never a leading digit.
func sanitizeAppName(name string) string {
	// Replace invalid characters with underscores
	name = invalidAppNameChars.ReplaceAllString(name, "_")
	// Remove consecutive underscores
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}
	// Trim leading/trailing underscores
	name = strings.Trim(name, "_")
	if name == "" {
		name = "wailsapp"
	}
	// Prefix with underscore if starts with digit. This has to happen after the
	// trim, which would otherwise strip the prefix again and leave an element
	// GTK refuses, e.g. "1Password" -> "org.wails.1password".
	name = leadingDigits.ReplaceAllString(name, "_$0")
	return strings.ToLower(name)
}

// maxApplicationIDLength is the longest id GTK accepts, inherited from the
// D-Bus bus name limit.
const maxApplicationIDLength = 255

// validateApplicationID returns an error describing why GTK would refuse id,
// following the same contract as g_application_id_is_valid():
//
//   - the id is composed of two or more elements separated by a '.', and every
//     element holds at least one character;
//   - every element contains only the ASCII characters A-Z, a-z, 0-9, '_' and
//     '-', and does not begin with a digit;
//   - the id is at most 255 characters long.
//
// GTK only asserts on this, so an invalid id makes gtk_application_new() return
// NULL and takes the process down later, far away from the option that caused it.
//
// See: https://docs.gtk.org/gio/type_func.Application.id_is_valid.html
func validateApplicationID(id string) error {
	if id == "" {
		return errors.New("application id is empty")
	}
	if len(id) > maxApplicationIDLength {
		return fmt.Errorf("application id %q is %d characters long, the maximum is %d", id, len(id), maxApplicationIDLength)
	}

	elements := strings.Split(id, ".")
	if len(elements) < 2 {
		return fmt.Errorf("application id %q needs at least two elements separated by a '.', for example \"com.example.MyApp\"", id)
	}

	for _, element := range elements {
		if element == "" {
			return fmt.Errorf("application id %q has an empty element: it must not start or end with a '.', or contain \"..\"", id)
		}
		if element[0] >= '0' && element[0] <= '9' {
			return fmt.Errorf("application id %q has the element %q starting with a digit", id, element)
		}
		for i := 0; i < len(element); i++ {
			if !isApplicationIDChar(element[i]) {
				return fmt.Errorf("application id %q contains the invalid character %q: only A-Z, a-z, 0-9, '_' and '-' are allowed", id, rune(element[i]))
			}
		}
	}

	return nil
}

func isApplicationIDChar(c byte) bool {
	return c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' ||
		c == '_' || c == '-'
}

// applicationID returns the id to build the GtkApplication with. Options.Linux
// wins when it sets one, so sandboxed builds can match the id their runtime
// expects; everything else keeps the derived "org.wails.<name>".
//
// An id GTK would reject is reported as an error together with the derived id,
// so callers can carry on with an id that works instead of crashing inside GTK.
func applicationID(options Options) (string, error) {
	derived := "org.wails." + sanitizeAppName(options.Name)
	if len(derived) > maxApplicationIDLength {
		// sanitizeAppName never emits a '.', so cutting the tail can only leave
		// characters that are legal in the middle of an element.
		derived = derived[:maxApplicationIDLength]
	}

	id := options.Linux.ApplicationID
	if id == "" {
		return derived, nil
	}
	if err := validateApplicationID(id); err != nil {
		return derived, err
	}
	return id, nil
}

// programName returns the name to hand to g_set_prgname, or "" to leave the
// program name at whatever GTK picked up from the executable.
//
// GTK takes the Wayland surface app_id from g_get_prgname(), so a window is only
// matched with its .desktop file when the program name carries the application
// id as well. An application that sets Options.Linux.ApplicationID inherits it
// here rather than having to repeat the same string in ProgramName.
//
// What is inherited is appID, the id the GtkApplication was built with, so the
// program name cannot disagree with it: an ApplicationID that validation
// rejected falls back to the derived id in both places. Without an
// ApplicationID nothing is derived, keeping the program name of applications
// that set neither option as it was.
func programName(options Options, appID string) string {
	if options.Linux.ProgramName != "" {
		return options.Linux.ProgramName
	}
	if options.Linux.ApplicationID != "" {
		return appID
	}
	return ""
}
