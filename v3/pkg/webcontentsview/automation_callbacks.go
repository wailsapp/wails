package webcontentsview

import (
	"encoding/json"
	"sync"
	"sync/atomic"
)

type automationCommandResult struct {
	payload string
	err     string
}

var automationCommandCallbacks sync.Map
var automationCommandCallbackID uintptr

func registerAutomationCommandCallback(ch chan automationCommandResult) uintptr {
	id := atomic.AddUintptr(&automationCommandCallbackID, 1)
	automationCommandCallbacks.Store(id, ch)
	return id
}

func dispatchAutomationCommandResult(id uintptr, result automationCommandResult) {
	if ch, ok := automationCommandCallbacks.Load(id); ok {
		ch.(chan automationCommandResult) <- result
		automationCommandCallbacks.Delete(id)
	}
}

func dispatchAutomationEvent(viewID uint, method string, payload string) {
	view := lookupAutomationView(viewID)
	if view == nil {
		return
	}

	switch method {
	case "Console.messageAdded":
		var message AutomationConsoleMessage
		if payload == "" || json.Unmarshal([]byte(payload), &message) != nil {
			return
		}
		view.recordAutomationConsoleMessage(message)

	case "Runtime.exceptionThrown":
		var exception AutomationException
		if payload == "" || json.Unmarshal([]byte(payload), &exception) != nil {
			return
		}
		view.recordAutomationException(exception)

	case "Page.domContentEventFired":
		var event struct {
			URL       string `json:"url"`
			Timestamp int64  `json:"timestamp"`
		}
		if payload == "" || json.Unmarshal([]byte(payload), &event) != nil {
			return
		}
		view.recordAutomationDOMContentLoaded(event.URL, event.Timestamp)

	case "Page.frameStartedLoading":
		var event struct {
			URL       string `json:"url"`
			Timestamp int64  `json:"timestamp"`
		}
		if payload == "" || json.Unmarshal([]byte(payload), &event) != nil {
			return
		}
		view.recordAutomationNavigationStarted(event.URL, event.Timestamp)

	case "Page.frameNavigated":
		var event struct {
			URL       string `json:"url"`
			Title     string `json:"title"`
			Timestamp int64  `json:"timestamp"`
		}
		if payload == "" || json.Unmarshal([]byte(payload), &event) != nil {
			return
		}
		view.recordAutomationNavigationCommitted(event.URL, event.Title, event.Timestamp)

	case "Page.loadEventFired":
		var event struct {
			URL       string `json:"url"`
			Title     string `json:"title"`
			Timestamp int64  `json:"timestamp"`
		}
		if payload == "" || json.Unmarshal([]byte(payload), &event) != nil {
			return
		}
		view.recordAutomationLoadFinished(event.URL, event.Title, event.Timestamp)

	case "Page.windowOpenRequested":
		var event struct {
			URL       string `json:"url"`
			Timestamp int64  `json:"timestamp"`
		}
		if payload == "" || json.Unmarshal([]byte(payload), &event) != nil {
			return
		}
		view.recordAutomationWindowOpenRequested(event.URL, event.Timestamp)

	case "Page.webContentProcessTerminated":
		var event struct {
			Timestamp int64 `json:"timestamp"`
		}
		if payload == "" || json.Unmarshal([]byte(payload), &event) != nil {
			return
		}
		view.recordAutomationProcessTerminated(event.Timestamp)

	case "Target.targetInfoChanged":
		var event struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		}
		if payload == "" || json.Unmarshal([]byte(payload), &event) != nil {
			return
		}
		view.recordAutomationPageState(event.URL, event.Title)
	}
}
