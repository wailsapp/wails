package main

import (
	"bytes"
	"os"
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

extern void systemEventHandler(char*);

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

	// Loop over each line in the file
	for _, line := range bytes.Split(eventNames, []byte{'\n'}) {

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// split on the colon
		split := bytes.Split(line, []byte{':'})
		platform := strings.TrimSpace(string(split[0]))
		event := strings.TrimSpace(string(split[1]))

		// Title case the event name
		eventTitle := string(bytes.ToUpper([]byte{event[0]})) + event[1:]

		// Add to buffer
		switch platform {
		case "mac":
			macEventsDecl.WriteString("\t" + eventTitle + " string\n")
			macEventsValues.WriteString("\t\t" + event + ": \"" + platform + ":" + event + "\",\n")
			cHeaderEvents.WriteString("#define Event" + eventTitle + " \"" + platform + ":" + event + "\"\n")
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

}
