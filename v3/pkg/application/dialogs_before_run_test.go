package application

import (
	"sync"
	"testing"
	"time"
)

// A dialog asked for before Run cannot be shown: main-thread work is queued
// until Run replays it, and the goroutine asking is usually the one that would
// have called Run, so waiting for that work never returns. These bound the
// wait rather than asserting on it - on the unfixed code the call parks
// forever, so a test that simply called it would hang the package.

const dialogReturnWait = 5 * time.Second

type warningRecorder struct {
	mu       sync.Mutex
	warnings []string
}

func (w *warningRecorder) record(msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.warnings = append(w.warnings, msg)
}

func (w *warningRecorder) count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.warnings)
}

// newApplicationBeforeRun installs an application that has not started: no
// platform implementation, which is exactly the state App.Run publishes one
// out of.
func newApplicationBeforeRun(t *testing.T) *warningRecorder {
	t.Helper()

	recorder := &warningRecorder{}
	previous := globalApplication
	globalApplication = &App{options: Options{WarningHandler: recorder.record}}
	t.Cleanup(func() { globalApplication = previous })

	return recorder
}

func TestDialogsBeforeRunReturnInsteadOfWaiting(t *testing.T) {
	tests := []struct {
		name string
		// call reports a wrong result rather than asserting, so nothing in the
		// goroutine touches *testing.T: on unfixed code it is still parked
		// when the test gives up on it.
		call func() string
	}{
		{
			name: "MessageDialog.Show",
			call: func() string {
				newMessageDialog(InfoDialogType).SetMessage("before Run").Show()
				return ""
			},
		},
		{
			name: "OpenFileDialog.PromptForSingleSelection",
			call: func() string {
				selection, err := new(OpenFileDialogStruct).PromptForSingleSelection()
				if selection != "" || err != nil {
					return "got (" + selection + ", " + errText(err) + "), want the zero values"
				}
				return ""
			},
		},
		{
			name: "OpenFileDialog.PromptForMultipleSelection",
			call: func() string {
				selections, err := new(OpenFileDialogStruct).PromptForMultipleSelection()
				if selections != nil || err != nil {
					return "got a result or an error, want the zero values"
				}
				return ""
			},
		},
		{
			name: "SaveFileDialog.PromptForSingleSelection",
			call: func() string {
				selection, err := new(SaveFileDialogStruct).PromptForSingleSelection()
				if selection != "" || err != nil {
					return "got (" + selection + ", " + errText(err) + "), want the zero values"
				}
				return ""
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := newApplicationBeforeRun(t)

			returned := make(chan string, 1)
			go func() { returned <- test.call() }()

			select {
			case problem := <-returned:
				if problem != "" {
					t.Error(problem)
				}
			case <-time.After(dialogReturnWait):
				t.Fatalf("%s did not return within %s", test.name, dialogReturnWait)
			}

			if recorder.count() != 1 {
				t.Errorf("recorded %d warnings, want 1 saying no dialog was shown", recorder.count())
			}
		})
	}
}

// With no application at all there is nothing to warn through, and the call
// still has to come back.
func TestDialogsWithNoApplicationReturn(t *testing.T) {
	previous := globalApplication
	globalApplication = nil
	t.Cleanup(func() { globalApplication = previous })

	returned := make(chan struct{})
	go func() {
		newMessageDialog(InfoDialogType).Show()
		close(returned)
	}()

	select {
	case <-returned:
	case <-time.After(dialogReturnWait):
		t.Fatalf("MessageDialog.Show did not return within %s", dialogReturnWait)
	}
}

func errText(err error) string {
	if err == nil {
		return "<nil>"
	}
	return err.Error()
}
