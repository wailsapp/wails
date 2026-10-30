package term

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pterm/pterm"
)

func TestLiteralAndFormattedMessages(t *testing.T) {
	originalInfo, originalWarning := pterm.Info, pterm.Warning
	originalError, originalSuccess := pterm.Error, pterm.Success
	t.Cleanup(func() {
		pterm.Info, pterm.Warning = originalInfo, originalWarning
		pterm.Error, pterm.Success = originalError, originalSuccess
	})
	var output bytes.Buffer
	pterm.Info.Writer, pterm.Warning.Writer = &output, &output
	pterm.Error.Writer, pterm.Success.Writer = &output, &output
	message := "100% done: https://example.com/a%20b literal %s %%"
	for name, print := range map[string]func(any){
		"Info": Info, "Warning": Warning, "Error": Error, "Success": Success,
	} {
		t.Run(name, func(t *testing.T) {
			output.Reset()
			print(message)
			if got := output.String(); !strings.Contains(got, message) {
				t.Errorf("message was changed: %q", got)
			}
		})
	}
	for name, print := range map[string]func(any, ...any){"Infof": Infof, "Warningf": Warningf} {
		t.Run(name, func(t *testing.T) {
			output.Reset()
			print("progress %d%%", 42)
			if got := output.String(); !strings.Contains(got, "progress 42%") {
				t.Errorf("formatted message = %q", got)
			}
			output.Reset()
			print("literal %%s")
			if got := output.String(); !strings.Contains(got, "literal %s") {
				t.Errorf("zero-argument format = %q", got)
			}
		})
	}
}
