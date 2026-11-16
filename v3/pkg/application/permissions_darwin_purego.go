//go:build darwin && purego && !ios && !server

package application

import "github.com/ebitengine/purego/objc"

// The purego counterpart of resolveMediaCapturePermission in
// permissions_darwin.go. The decision itself (mediaCaptureDecision and the
// rest of permissions_darwin_capture.go) is shared; this file only maps
// WebKit's WKMediaCaptureType onto it and answers the decision handler.

// wkMediaCaptureType mirrors WKMediaCaptureType as WebKit passes it to
// requestMediaCapturePermissionForOrigin:...type:decisionHandler:.
type wkMediaCaptureType int

const (
	wkMediaCaptureTypeCamera              wkMediaCaptureType = 0
	wkMediaCaptureTypeMicrophone          wkMediaCaptureType = 1
	wkMediaCaptureTypeCameraAndMicrophone wkMediaCaptureType = 2
)

// mediaCaptureNeeds returns which devices a getUserMedia request asks for. It
// is the same mapping the cgo WKUIDelegate performs in
// webview_window_darwin.m. A type it does not know asks for neither device,
// which mediaCaptureDecision answers with WebKit's prompt.
func mediaCaptureNeeds(mediaType int) (needAudio bool, needVideo bool) {
	needAudio = mediaType == int(wkMediaCaptureTypeMicrophone) || mediaType == int(wkMediaCaptureTypeCameraAndMicrophone)
	needVideo = mediaType == int(wkMediaCaptureTypeCamera) || mediaType == int(wkMediaCaptureTypeCameraAndMicrophone)
	return needAudio, needVideo
}

// mediaCaptureDecisionForType is the whole decision path for one request: the
// window's Permissions, filtered to the devices the request asks for, reduced
// to the WKPermissionDecision value the handler receives.
func mediaCaptureDecisionForType(windowID uint, mediaType int) captureDecision {
	needAudio, needVideo := mediaCaptureNeeds(mediaType)
	return mediaCaptureDecision(windowID, needAudio, needVideo)
}

// webViewRequestMediaCapturePermission implements
// -[WKUIDelegate webView:requestMediaCapturePermissionForOrigin:initiatedByFrame:type:decisionHandler:].
//
// This selector is macOS 12+. On older systems WebKit never calls it and
// falls back to its own prompt, so registering it unconditionally is safe;
// the guard is only needed where a selector is called, not where it is
// implemented.
//
// The decision handler is invoked synchronously, before this method returns,
// as WebKit requires; it is not retained or copied.
func webViewRequestMediaCapturePermission(self objc.ID, cmd objc.SEL, wv objc.ID, origin objc.ID, frame objc.ID, mediaType int, decisionHandler objc.ID) {
	decision := captureDecisionPrompt
	// The delegate's window id is the purego equivalent of the cgo
	// windowId ivar. A delegate with no known window keeps WebKit's prompt,
	// as an unknown window does in resolvePermission.
	if windowID, ok := windowIDForDelegate(self); ok {
		decision = mediaCaptureDecisionForType(windowID, mediaType)
	}
	invokeForeignBlock(decisionHandler, uintptr(decision))
}

// mediaCapturePermissionMethod is the MethodDef for the selector above, added
// to the window delegate class (which is also the window's WKUIDelegate).
func mediaCapturePermissionMethod() objc.MethodDef {
	return objc.MethodDef{
		Cmd: sel_("webView:requestMediaCapturePermissionForOrigin:initiatedByFrame:type:decisionHandler:"),
		Fn:  webViewRequestMediaCapturePermission,
	}
}
