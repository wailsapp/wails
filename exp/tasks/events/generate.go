package main

import (
	"bytes"
	"os"
	"strconv"
	"strings"
)

var eventsGo = `package events

var Mac = newMacEvents()

type macEvents struct {
$$MACEVENTSDECL}

func newMacEvents() macEvents {
	return macEvents{
$$MACEVENTSVALUES	}
}
`

var eventsH = `//go:build darwin

#ifndef _events_h
#define _events_h

extern void applicationEventHandler(unsigned int);
extern void windowEventHandler(unsigned int, unsigned int);

$$CHEADEREVENTS

#endif`

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

	// Loop over each line in the file
	for id, line := range bytes.Split(eventNames, []byte{'\n'}) {

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
			macEventsDecl.WriteString("\t" + eventTitle + " uint\n")
			macEventsValues.WriteString("\t\t" + event + ": " + strconv.Itoa(id) + ",\n")
			cHeaderEvents.WriteString("#define Event" + eventTitle + " " + strconv.Itoa(id) + "\n")
			if ignoreEvent {
				continue
			}
			// Check if this is a window event
			if strings.HasPrefix(event, "Window") {
				windowDelegateEvents.WriteString(`- (void)` + delegateEventFunction + `:(NSNotification *)notification {
    windowEventHandler(self.windowId, Event` + eventTitle + `);
}

`)
			}
			if strings.HasPrefix(event, "Application") {
				applicationDelegateEvents.WriteString(`- (void)` + delegateEventFunction + `:(NSNotification *)notification {
    applicationEventHandler(Event` + eventTitle + `);
}

`)
			}
		}

	}

	// Save the eventsGo template substituting the values and decls
	templateToWrite := strings.ReplaceAll(eventsGo, "$$MACEVENTSDECL", macEventsDecl.String())
	templateToWrite = strings.ReplaceAll(templateToWrite, "$$MACEVENTSVALUES", macEventsValues.String())
	err = os.WriteFile("../../pkg/events/events.go", []byte(templateToWrite), 0644)
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
	windowDelegate, err := os.ReadFile("../../pkg/application/window_delegate.m")
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
			continue
		}
		if bytes.Contains(line, []byte("// GENERATED EVENTS END")) {
			inGeneratedEvents = false
			buffer.WriteString("// GENERATED EVENTS END\n")
			continue
		}
		if !inGeneratedEvents {
			buffer.Write(line)
			buffer.WriteString("\n")
		}
	}
	err = os.WriteFile("../../pkg/application/window_delegate.m", buffer.Bytes(), 0755)
	if err != nil {
		panic(err)
	}

	// Load the app_delegate.m file
	appDelegate, err := os.ReadFile("../../pkg/application/app_delegate.m")
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
			buffer.Write(line)
			buffer.WriteString("\n")
		}
	}
	err = os.WriteFile("../../pkg/application/app_delegate.m", buffer.Bytes(), 0755)
	if err != nil {
		panic(err)
	}

}
