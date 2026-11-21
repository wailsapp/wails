//go:build openharmony && !server

package application

import "encoding/json"

// The clipboard is backed by the system pasteboard via the WailsBridge.

type openharmonyClipboardImpl struct{}

func newClipboardImpl() clipboardImpl {
	return &openharmonyClipboardImpl{}
}

func (c *openharmonyClipboardImpl) setText(text string) bool {
	payload, _ := json.Marshal(map[string]string{"text": text})
	ohBridgeVoid("setClipboardText", string(payload))
	return true
}

func (c *openharmonyClipboardImpl) text() (string, bool) {
	res, ok := ohBridgeString("getClipboardText", "{}")
	if !ok || res == "" {
		return "", false
	}
	return res, true
}
