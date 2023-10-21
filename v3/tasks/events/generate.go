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

var eventsH = `//go:build darwin

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
	cHeaderEvents := bytes.NewBufferString("")
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

	eventToJS := bytes.NewBufferString("")

	var id int
	var maxMacEvents int
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
			event = event[:len(event)-1]
			ignoreEvent = true
		}

		// Title case the event name
		eventTitle := string(bytes.ToUpper([]byte{event[0]})) + event[1:]
		// delegate function name has a lowercase first character
		delegateEventFunction := string(bytes.ToLower([]byte{event[0]})) + event[1:]

		// Add to buffer
		switch platform {
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
			cHeaderEvents.WriteString("#define Event" + eventTitle + " " + strconv.Itoa(id) + "\n")
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
			eventToJS.WriteString("\t" + strconv.Itoa(id) + ": \"" + strings.TrimSpace(string(line)) + "\",\n")
		}
	}

	cHeaderEvents.WriteString("\n#define MAX_EVENTS " + strconv.Itoa(maxMacEvents+1) + "\n")

	// Save the eventsGo template substituting the values and decls
	templateToWrite := strings.ReplaceAll(eventsGo, "$$MACEVENTSDECL", macEventsDecl.String())
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
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$COMMONJSEVENTS", commonJSEvents.String())
	err = os.WriteFile("../../internal/runtime/desktop/api/event_types.js", []byte(templateToWrite), 0644)
	if err != nil {
		panic(err)
	}

	// Save the eventsH template substituting the values and decls
	templateToWrite = strings.ReplaceAll(eventsH, "$$CHEADEREVENTS", cHeaderEvents.String())
	err = os.WriteFile("../../pkg/events/events.h", []byte(templateToWrite), 0644)
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
