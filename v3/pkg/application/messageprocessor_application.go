package application

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ApplicationHide = 0
	ApplicationShow = 1
	ApplicationQuit = 2

	// This should match ApplicationFilesDroppedWithContext in system.ts
	applicationMethodFilesDroppedWithContext = 100
)

// Define a struct for the JSON payload from HandlePlatformFileDrop
type fileDropPayload struct {
	Filenames []string `json:"filenames"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	ElementID string   `json:"elementId"`
	ClassList []string `json:"classList"`
}

var applicationMethodNames = map[int]string{
	ApplicationQuit: "Quit",
	ApplicationHide: "Hide",
	ApplicationShow: "Show",
	applicationMethodFilesDroppedWithContext: "FilesDroppedWithContext", // For logging
}

func (m *MessageProcessor) processApplicationMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	switch method {
	case ApplicationQuit:
		globalApplication.Quit()
		m.ok(rw)
	case ApplicationHide:
		globalApplication.Hide()
		m.ok(rw)
	case ApplicationShow:
		globalApplication.Show()
		m.ok(rw)
	case applicationMethodFilesDroppedWithContext:
		m.Info("[DragDropDebug] processApplicationMethod: Entered applicationMethodFilesDroppedWithContext case")
		var payload fileDropPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			m.httpError(rw, "Error decoding file drop payload:", err)
			return
		}
		m.Info("[DragDropDebug] processApplicationMethod: Decoded payload", "payload", fmt.Sprintf("%+v", payload))

		dropDetails := &DropZoneDetails{
			X:         payload.X,
			Y:         payload.Y,
			ElementID: payload.ElementID,
			ClassList: payload.ClassList,
		}

		wvWindow, ok := window.(*WebviewWindow)
		if !ok {
			m.httpError(rw, "Error: Target window is not a WebviewWindow for FilesDroppedWithContext", nil)
			return
		}

		msg := &dragAndDropMessage{
			windowId:  wvWindow.id,
			filenames: payload.Filenames,
			DropZone:  dropDetails,
		}

		m.Info("[DragDropDebug] processApplicationMethod: Sending message to windowDragAndDropBuffer", "message", fmt.Sprintf("%+v", msg))
		windowDragAndDropBuffer <- msg
		m.ok(rw)
	default:
		m.httpError(rw, "Invalid application call:", fmt.Errorf("unknown method: %d", method))
		return
	}

	m.Info("Runtime call:", "method", "Application."+applicationMethodNames[method])
}
