package main

import (
	"bytes"
	"github.com/Masterminds/semver/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
	"strconv"
	"strings"
)

var eventsGo = `package events

type ApplicationEventType uint
type WindowEventType      uint

var Common = newCommonEvents()

type commonEvents struct {
$$COMMONEVENTSDECL}

func newCommonEvents() commonEvents {
	return commonEvents{
$$COMMONEVENTSVALUES	}
}

var Linux = newLinuxEvents()

type linuxEvents struct {
$$LINUXEVENTSDECL}

func newLinuxEvents() linuxEvents {
	return linuxEvents{
$$LINUXEVENTSVALUES	}
}

var Mac = newMacEvents()

type macEvents struct {
$$MACEVENTSDECL}

func newMacEvents() macEvents {
	return macEvents{
$$MACEVENTSVALUES	}
}

var Windows = newWindowsEvents()

type windowsEvents struct {
$$WINDOWSEVENTSDECL}

func newWindowsEvents() windowsEvents {
	return windowsEvents{
$$WINDOWSEVENTSVALUES	}
}

func JSEvent(event uint) string {
	return eventToJS[event]
}

var eventToJS = map[uint]string{
$$EVENTTOJS}

`

var darwinEventsH = `//go:build darwin

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

$$CHEADEREVENTS

#endif`

var linuxEventsH = `//go:build linux

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

$$CHEADEREVENTS

#endif`

var eventsJS = `
export const EventTypes = {
	Windows: {
$$WINDOWSJSEVENTS	},
	Mac: {
$$MACJSEVENTS	},
	Linux: {
$$LINUXJSEVENTS	},
	Common: {
$$COMMONJSEVENTS	},
};
`

var eventsTS = `
export declare const EventTypes: {
	Windows: {
$$WINDOWSTSEVENTS	},
	Mac: {
$$MACTSEVENTS	},
	Linux: {
$$LINUXTSEVENTS	},
	Common: {
$$COMMONTSEVENTS	},
};
`

func main() {

	eventNames, err := os.ReadFile("../../pkg/events/events.txt")
	if err != nil {
		panic(err)
	}

	linuxEventsDecl := bytes.NewBufferString("")
	linuxEventsValues := bytes.NewBufferString("")
	linuxCHeaderEvents := bytes.NewBufferString("")

	macEventsDecl := bytes.NewBufferString("")
	macEventsValues := bytes.NewBufferString("")
	macCHeaderEvents := bytes.NewBufferString("")
	windowDelegateEvents := bytes.NewBufferString("")
	applicationDelegateEvents := bytes.NewBufferString("")
	webviewDelegateEvents := bytes.NewBufferString("")

	windowsEventsDecl := bytes.NewBufferString("")
	windowsEventsValues := bytes.NewBufferString("")

	commonEventsDecl := bytes.NewBufferString("")
	commonEventsValues := bytes.NewBufferString("")

	linuxJSEvents := bytes.NewBufferString("")
	macJSEvents := bytes.NewBufferString("")
	windowsJSEvents := bytes.NewBufferString("")
	commonJSEvents := bytes.NewBufferString("")

	linuxTSEvents := bytes.NewBufferString("")
	macTSEvents := bytes.NewBufferString("")
	windowsTSEvents := bytes.NewBufferString("")
	commonTSEvents := bytes.NewBufferString("")

	eventToJS := bytes.NewBufferString("")

	var id int
	//	var maxLinuxEvents int
	var maxMacEvents int
	var maxLinuxEvents int
	var line []byte
	// Loop over each line in the file
	for id, line = range bytes.Split(eventNames, []byte{'\n'}) {

		// First 1024 is reserved
		id = id + 1024

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// split on the colon
		split := bytes.Split(line, []byte{':'})
		platform := strings.TrimSpace(string(split[0]))
		event := strings.TrimSpace(string(split[1]))
		var ignoreEvent bool
		if strings.HasSuffix(event, "!") {
			event = strings.TrimSuffix(event, "!")
			ignoreEvent = true
		}
		// Strip last byte of line if it's a "!" character
		if line[len(line)-1] == '!' {
			line = line[:len(line)-1]
		}

		// Title case the event name
		eventTitle := string(bytes.ToUpper([]byte{event[0]})) + event[1:]
		// delegate function name has a lowercase first character
		delegateEventFunction := string(bytes.ToLower([]byte{event[0]})) + event[1:]

		// Add to buffer
		switch platform {
		case "linux":
			eventType := "ApplicationEventType"
			if strings.HasPrefix(event, "Window") {
				eventType = "WindowEventType"
			}
			if strings.HasPrefix(event, "WebView") {
				eventType = "WindowEventType"
			}
			linuxEventsDecl.WriteString("\t" + eventTitle + " " + eventType + "\n")
			linuxEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(id) + ",\n")
			linuxJSEvents.WriteString("\t\t" + event + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			linuxTSEvents.WriteString("\t\t" + event + ": string,\n")
			eventToJS.WriteString("\t" + strconv.Itoa(id) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			maxLinuxEvents = id
			linuxCHeaderEvents.WriteString("#define Event" + eventTitle + " " + strconv.Itoa(id) + "\n")
		case "mac":
			eventType := "ApplicationEventType"
			if strings.HasPrefix(event, "Window") {
				eventType = "WindowEventType"
			}
			if strings.HasPrefix(event, "WebView") {
				eventType = "WindowEventType"
			}
			macEventsDecl.WriteString("\t" + eventTitle + " " + eventType + "\n")
			macEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(id) + ",\n")
			macJSEvents.WriteString("\t\t" + event + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			macTSEvents.WriteString("\t\t" + event + ": string,\n")
			macCHeaderEvents.WriteString("#define Event" + eventTitle + " " + strconv.Itoa(id) + "\n")
			eventToJS.WriteString("\t" + strconv.Itoa(id) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			maxMacEvents = id
			if ignoreEvent {
				continue
			}
			// Check if this is a window event
			if strings.HasPrefix(event, "Window") {
				windowDelegateEvents.WriteString(`- (void)` + delegateEventFunction + `:(NSNotification *)notification {
    if( hasListeners(Event` + eventTitle + `) ) {
        processWindowEvent(self.windowId, Event` + eventTitle + `);
    }
}

`)
			}
			// Check if this is a webview event
			if strings.HasPrefix(event, "WebView") {
				webViewFunction := strings.TrimPrefix(event, "WebView")
				webViewFunction = string(bytes.ToLower([]byte{webViewFunction[0]})) + webViewFunction[1:]
				webviewDelegateEvents.WriteString(`- (void)webView:(WKWebView *)webview ` + webViewFunction + `:(WKNavigation *)navigation {
    if( hasListeners(Event` + eventTitle + `) ) {
        processWindowEvent(self.windowId, Event` + eventTitle + `);
    }
}

`)
			}
			if strings.HasPrefix(event, "Application") {
				applicationDelegateEvents.WriteString(`- (void)` + delegateEventFunction + `:(NSNotification *)notification {
    if( hasListeners(Event` + eventTitle + `) ) {
        processApplicationEvent(Event` + eventTitle + `, NULL);
    }
}

`)
			}
		case "common":
			eventType := "ApplicationEventType"
			if strings.HasPrefix(event, "Window") {
				eventType = "WindowEventType"
			}
			if strings.HasPrefix(event, "WebView") {
				eventType = "WindowEventType"
			}
			commonEventsDecl.WriteString("\t" + eventTitle + " " + eventType + "\n")
			commonEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(id) + ",\n")
			commonJSEvents.WriteString("\t\t" + event + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			commonTSEvents.WriteString("\t\t" + event + ": string,\n")
			eventToJS.WriteString("\t" + strconv.Itoa(id) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
		case "windows":
			eventType := "ApplicationEventType"
			if strings.HasPrefix(event, "Window") {
				eventType = "WindowEventType"
			}
			if strings.HasPrefix(event, "WebView") {
				eventType = "WindowEventType"
			}
			windowsEventsDecl.WriteString("\t" + eventTitle + " " + eventType + "\n")
			windowsEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(id) + ",\n")
			windowsJSEvents.WriteString("\t\t" + event + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			windowsTSEvents.WriteString("\t\t" + event + ": string,\n")
			eventToJS.WriteString("\t" + strconv.Itoa(id) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
		}
	}

	macCHeaderEvents.WriteString("\n#define MAX_EVENTS " + strconv.Itoa(maxMacEvents+1) + "\n")
	linuxCHeaderEvents.WriteString("\n#define MAX_EVENTS " + strconv.Itoa(maxLinuxEvents+1) + "\n")

	// Save the eventsGo template substituting the values and decls
	templateToWrite := strings.ReplaceAll(eventsGo, "$$LINUXEVENTSDECL", linuxEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$LINUXEVENTSVALUES", linuxEventsValues.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$MACEVENTSDECL", macEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$MACEVENTSVALUES", macEventsValues.String())

	templateToWrite = strings.ReplaceAll(templateToWrite, "$$WINDOWSEVENTSDECL", windowsEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$WINDOWSEVENTSVALUES", windowsEventsValues.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONEVENTSDECL", commonEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONEVENTSVALUES", commonEventsValues.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$EVENTTOJS", eventToJS.String())
	err = os.WriteFile("../../pkg/events/events.go", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Save the eventsJS template substituting the values and decls
	templateToWrite = strings.ReplaceAll(eventsJS, "$$MACJSEVENTS", macJSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$WINDOWSJSEVENTS", windowsJSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$LINUXJSEVENTS", linuxJSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONJSEVENTS", commonJSEvents.String())
	err = os.WriteFile("../../internal/runtime/desktop/@wailsio/runtime/src/event_types.js", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Save the eventsTS template substituting the values and decls
	templateToWrite = strings.ReplaceAll(eventsTS, "$$MACTSEVENTS", macTSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$WINDOWSTSEVENTS", windowsTSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$LINUXTSEVENTS", linuxTSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONTSEVENTS", commonTSEvents.String())
	err = os.WriteFile("../../internal/runtime/desktop/@wailsio/runtime/types/event_types.d.ts", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Save the darwinEventsH template substituting the values and decls
	templateToWrite = strings.ReplaceAll(darwinEventsH, "$$CHEADEREVENTS", macCHeaderEvents.String())
	err = os.WriteFile("../../pkg/events/events_darwin.h", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Save the linuxEventsH template substituting the values and decls
	templateToWrite = strings.ReplaceAll(linuxEventsH, "$$CHEADEREVENTS", linuxCHeaderEvents.String())
	err = os.WriteFile("../../pkg/events/events_linux.h", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Load the window_delegate.m file
	windowDelegate, err := os.ReadFile("../../pkg/application/webview_window_darwin.m")
	if err != nil {
		panic(err)
	}
	// iterate over the lines until we reach a line that says "// GENERATED EVENTS START"
	// then we insert the events
	// then we iterate until we reach a line that says "// GENERATED EVENTS END"
	// then we write the file
	var buffer bytes.Buffer
	var inGeneratedEvents bool
	for _, line := range bytes.Split(windowDelegate, []byte{'\n'}) {
		if bytes.Contains(line, []byte("// GENERATED EVENTS START")) {
			inGeneratedEvents = true
			buffer.WriteString("// GENERATED EVENTS START\n")
			buffer.WriteString(windowDelegateEvents.String())
			buffer.WriteString(webviewDelegateEvents.String())
			continue
		}
		if bytes.Contains(line, []byte("// GENERATED EVENTS END")) {
			inGeneratedEvents = false
			buffer.WriteString("// GENERATED EVENTS END\n")
			continue
		}
		if !inGeneratedEvents {
			if len(line) > 0 {
				buffer.Write(line)
				buffer.WriteString("\n")
			}
		}
	}
	err = os.WriteFile("../../pkg/application/webview_window_darwin.m", buffer.Bytes(), 0755)
	if err != nil {
		panic(err)
	}

	// Load the app_delegate.m file
	appDelegate, err := os.ReadFile("../../pkg/application/application_darwin_delegate.m")
	if err != nil {
		panic(err)
	}
	// iterate over the lines until we reach a line that says "// GENERATED EVENTS START"
	// then we insert the events
	// then we iterate until we reach a line that says "// GENERATED EVENTS END"
	// then we write the file
	buffer.Reset()
	for _, line := range bytes.Split(appDelegate, []byte{'\n'}) {
		if bytes.Contains(line, []byte("// GENERATED EVENTS START")) {
			inGeneratedEvents = true
			buffer.WriteString("// GENERATED EVENTS START\n")
			buffer.WriteString(applicationDelegateEvents.String())
			continue
		}
		if bytes.Contains(line, []byte("// GENERATED EVENTS END")) {
			inGeneratedEvents = false
			buffer.WriteString("// GENERATED EVENTS END\n")
			continue
		}
		if !inGeneratedEvents {
			if len(line) > 0 {
				buffer.Write(line)
				buffer.WriteString("\n")
			}
		}
	}
	err = os.WriteFile("../../pkg/application/application_darwin_delegate.m", buffer.Bytes(), 0755)
	if err != nil {
		panic(err)
	}

	// Load the runtime package.json
	packageJsonFilename := "../../internal/runtime/desktop/@wailsio/runtime/package.json"
	packageJSON, err := os.ReadFile(packageJsonFilename)
	if err != nil {
		panic(err)
	}
	version := gjson.Get(string(packageJSON), "version").String()
	// Parse and increment version
	v := semver.MustParse(version)
	prerelease := v.Prerelease()
	// Split the prerelease by the "." and increment the last part by 1
	parts := strings.Split(prerelease, ".")
	prereleaseDigits, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		panic(err)
	}
	prereleaseNumber := strconv.Itoa(prereleaseDigits + 1)
	parts[len(parts)-1] = prereleaseNumber
	prerelease = strings.Join(parts, ".")
	newVersion, err := v.SetPrerelease(prerelease)
	if err != nil {
		panic(err)
	}

	// Set new version using sjson
	newJSON, err := sjson.Set(string(packageJSON), "version", newVersion.String())
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(packageJsonFilename, []byte(newJSON), 0644)
}
