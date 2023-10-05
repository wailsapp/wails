package main

import (
	"bytes"
	"os"
	"strconv"
	"strings"
)

var eventsGo = `package events

type ApplicationEventType uint
type WindowEventType      uint

const (
	FilesDropped WindowEventType = iota
)

var Common = newCommonEvents()

type commonEvents struct {
$$COMMONEVENTSDECL}

func newCommonEvents() commonEvents {
	return commonEvents{
$$COMMONEVENTSVALUES	}
}

var Mac = newMacEvents()

type macEvents struct {
$$MACEVENTSDECL}

func newMacEvents() macEvents {
	return macEvents{
$$MACEVENTSVALUES	}
}

var Linux = newLinuxEvents()

type linuxEvents struct {
$$LINUXEVENTSDECL}

func newLinuxEvents() linuxEvents {
	return linuxEvents{
$$LINUXEVENTSVALUES	}
}

var Windows = newWindowsEvents()

type windowsEvents struct {
$$WINDOWSEVENTSDECL}

func newWindowsEvents() windowsEvents {
	return windowsEvents{
$$WINDOWSEVENTSVALUES	}
}

func JSEvent(platform string, event uint) string {
	return eventToJS[platform][event]
}

var eventToJS = map[string]map[uint]string{
	"windows": {
$$WINDOWSEVENTTOJS	},
	"mac": {
$$MACEVENTTOJS	},
	"linux": {
$$LINUXEVENTTOJS	},
	"common": {
$$COMMONEVENTTOJS	},
}

`

var eventsDarwinH = `//go:build darwin

#ifndef _events_darwin_h
#define _events_darwin_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

$$CDARWINHEADEREVENTS

#endif`

var eventsLinuxH = `//go:build linux

#ifndef _events_linux_h
#define _events_linux_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

$$CLINUXHEADEREVENTS

#endif`

var eventsJS = `
export const EventTypes = {
	Windows: {
$$WINDOWSJSEVENTS	},
	Mac: {
$$MACJSEVENTS	},
	Linux: {
$$LINUXJSEVENTS},
	Common: {
$$COMMONJSEVENTS	},
};
`

func main() {

	eventNames, err := os.ReadFile("../../pkg/events/events.txt")
	if err != nil {
		panic(err)
	}

	macEventsDecl := bytes.NewBufferString("")
	macEventsValues := bytes.NewBufferString("")
	cDarwinHeaderEvents := bytes.NewBufferString("")
	windowDelegateEvents := bytes.NewBufferString("")
	applicationDelegateEvents := bytes.NewBufferString("")
	webviewDelegateEvents := bytes.NewBufferString("")

	windowsEventsDecl := bytes.NewBufferString("")
	windowsEventsValues := bytes.NewBufferString("")

	commonEventsDecl := bytes.NewBufferString("")
	commonEventsValues := bytes.NewBufferString("")

	macJSEvents := bytes.NewBufferString("")
	windowsJSEvents := bytes.NewBufferString("")
	commonJSEvents := bytes.NewBufferString("")

	windowsEventToJS := bytes.NewBufferString("")
	macEventToJS := bytes.NewBufferString("")
	commonEventToJS := bytes.NewBufferString("")
	linuxEventToJS := bytes.NewBufferString("")

	linuxEventsDecl := bytes.NewBufferString("")
	linuxEventsValues := bytes.NewBufferString("")
	linuxJSEvents := bytes.NewBufferString("")
	cLinuxHeaderEvents := bytes.NewBufferString("")

	var currentMacEventNumber int
	var currentLinuxEventNumber int
	currentCommonEventNumber := 4096
	var currentWindowsEventNumber int
	var line []byte
	// Loop over each line in the file
	for _, line = range bytes.Split(eventNames, []byte{'\n'}) {

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
			event = event[:len(event)-1]
			ignoreEvent = true
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
			cLinuxHeaderEvents.WriteString("#define Event" + eventTitle + " " + strconv.Itoa(currentLinuxEventNumber) + "\n")
			linuxEventsDecl.WriteString("\t" + eventTitle + " " + eventType + "\n")
			linuxEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(currentLinuxEventNumber) + ",\n")
			linuxJSEvents.WriteString("\t\t" + event + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			if ignoreEvent {
				continue
			}
			linuxEventToJS.WriteString("\t\t" + strconv.Itoa(currentLinuxEventNumber) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			currentLinuxEventNumber++
		case "mac":
			eventType := "ApplicationEventType"
			if strings.HasPrefix(event, "Window") {
				eventType = "WindowEventType"
			}
			if strings.HasPrefix(event, "WebView") {
				eventType = "WindowEventType"
			}
			macEventsDecl.WriteString("\t" + eventTitle + " " + eventType + "\n")
			macEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(currentMacEventNumber) + ",\n")
			macJSEvents.WriteString("\t\t" + event + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			cDarwinHeaderEvents.WriteString("#define Event" + eventTitle + " " + strconv.Itoa(currentMacEventNumber) + "\n")
			if ignoreEvent {
				continue
			}
			macEventToJS.WriteString("\t\t" + strconv.Itoa(currentMacEventNumber) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			currentMacEventNumber++
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
			commonEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(currentCommonEventNumber) + ",\n")
			commonJSEvents.WriteString("\t\t" + event + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			commonEventToJS.WriteString("\t\t" + strconv.Itoa(currentCommonEventNumber) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			currentCommonEventNumber++
		case "windows":
			eventType := "ApplicationEventType"
			if strings.HasPrefix(event, "Window") {
				eventType = "WindowEventType"
			}
			if strings.HasPrefix(event, "WebView") {
				eventType = "WindowEventType"
			}
			windowsEventsDecl.WriteString("\t" + eventTitle + " " + eventType + "\n")
			windowsEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(currentWindowsEventNumber) + ",\n")
			windowsJSEvents.WriteString("\t\t" + event + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			windowsEventToJS.WriteString("\t\t" + strconv.Itoa(currentWindowsEventNumber) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
			currentWindowsEventNumber++
		}
	}

	cLinuxHeaderEvents.WriteString("\n#define MAX_EVENTS " + strconv.Itoa(currentLinuxEventNumber) + "\n")
	cDarwinHeaderEvents.WriteString("\n#define MAX_EVENTS " + strconv.Itoa(currentMacEventNumber) + "\n")

	// Save the eventsGo template substituting the values and decls
	templateToWrite := strings.ReplaceAll(eventsGo, "$$MACEVENTSDECL", macEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$LINUXEVENTSDECL", linuxEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$LINUXEVENTSVALUES", linuxEventsValues.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$MACEVENTSVALUES", macEventsValues.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$WINDOWSEVENTSDECL", windowsEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$WINDOWSEVENTSVALUES", windowsEventsValues.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONEVENTSDECL", commonEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONEVENTSVALUES", commonEventsValues.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$WINDOWSEVENTTOJS", windowsEventToJS.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$MACEVENTTOJS", macEventToJS.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$LINUXEVENTTOJS", linuxEventToJS.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONEVENTTOJS", commonEventToJS.String())
	err = os.WriteFile("../../pkg/events/events.go", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Save the eventsJS template substituting the values and decls
	templateToWrite = strings.ReplaceAll(eventsJS, "$$MACJSEVENTS", macJSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$WINDOWSJSEVENTS", windowsJSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$LINUXJSEVENTS", linuxJSEvents.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONJSEVENTS", commonJSEvents.String())
	err = os.WriteFile("../../internal/runtime/desktop/api/event_types.js", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Save the eventsDarwinH template substituting the values and decls
	templateToWrite = strings.ReplaceAll(eventsDarwinH, "$$CDARWINHEADEREVENTS", cDarwinHeaderEvents.String())
	err = os.WriteFile("../../pkg/events/events_darwin.h", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Save the eventsDarwinH template substituting the values and decls
	templateToWrite = strings.ReplaceAll(eventsLinuxH, "$$CLINUXHEADEREVENTS", cLinuxHeaderEvents.String())
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

}
