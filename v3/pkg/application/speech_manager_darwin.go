//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework AVFoundation -framework Speech

#include <stdlib.h>
#include "speech_manager_darwin.h"
*/
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"unsafe"
)

// speechEvent carries AVSpeechSynthesizer and SFSpeechRecognizer callbacks
// from the main thread to the Go drain loop registered in init.
type speechEvent struct {
	kind      speechEventKind
	id        uint64
	text      string
	cancelled bool
	final     bool
	err       error
}

type speechEventKind int

const (
	speechEventUtteranceFinished speechEventKind = iota
	speechEventRecognitionResult
	speechEventRecognitionError
)

var speechEvents = make(chan speechEvent, 64)

// speechAuthorizationEvents receives the outcome of an authorization prompt.
// Only one Recognize call runs at a time (speechRecognizeLock), so a single
// channel suffices.
type speechAuthorizationEvent struct {
	kind   int
	status int
}

var speechAuthorizationEvents = make(chan speechAuthorizationEvent, 4)

//export speechUtteranceFinishedCallback
func speechUtteranceFinishedCallback(id C.ulonglong, cancelled C.bool) {
	speechEvents <- speechEvent{kind: speechEventUtteranceFinished, id: uint64(id), cancelled: bool(cancelled)}
}

//export speechAuthorizationCallback
func speechAuthorizationCallback(kind C.int, status C.int) {
	select {
	case speechAuthorizationEvents <- speechAuthorizationEvent{kind: int(kind), status: int(status)}:
	default:
	}
}

//export speechRecognitionResultCallback
func speechRecognitionResultCallback(id C.ulonglong, text *C.char, final C.bool) {
	speechEvents <- speechEvent{kind: speechEventRecognitionResult, id: uint64(id), text: C.GoString(text), final: bool(final)}
}

//export speechRecognitionErrorCallback
func speechRecognitionErrorCallback(id C.ulonglong, message *C.char) {
	speechEvents <- speechEvent{kind: speechEventRecognitionError, id: uint64(id), err: errors.New(C.GoString(message))}
}

// handleSpeechEvent runs on the drain goroutine. Events are processed in
// order so partial transcripts never overtake each other.
func handleSpeechEvent(event speechEvent) {
	switch event.kind {
	case speechEventUtteranceFinished:
		if utterance := speechUtteranceByID(event.id); utterance != nil {
			speechUtterances.Delete(event.id)
			utterance.finish(event.cancelled)
		}
	case speechEventRecognitionResult:
		if session := speechRecognitionByID(event.id); session != nil {
			session.update(event.text, event.final)
		}
	case speechEventRecognitionError:
		if session := speechRecognitionByID(event.id); session != nil {
			session.end(session.Text(), event.err, true)
		}
	}
}

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			handleSpeechEvent(<-speechEvents)
		}
	})
}

func speechSynthesisSupported() bool {
	return bool(C.speechSynthesisAvailable())
}

func speechRecognitionSupported() bool {
	return bool(C.speechRecognitionAvailable())
}

func speechSpeak(id uint64, text string, opts SpeechOptions) error {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	cVoice := C.CString(opts.Voice)
	defer C.free(unsafe.Pointer(cVoice))
	C.speechSpeak(C.ulonglong(id), cText, cVoice, C.double(opts.Rate), C.double(opts.Volume))
	return nil
}

func speechStop(id uint64) {
	C.speechStopUtterance(C.ulonglong(id))
}

func speechStopAll() {
	C.speechStopAllUtterances()
}

func speechVoices() []Voice {
	raw := C.speechVoicesJSON()
	if raw == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(raw))
	var voices []Voice
	if err := json.Unmarshal([]byte(C.GoString(raw)), &voices); err != nil {
		globalApplication.error("speech: cannot decode voices: %v", err)
		return nil
	}
	sort.Slice(voices, func(i, j int) bool {
		if voices[i].Language != voices[j].Language {
			return voices[i].Language < voices[j].Language
		}
		return voices[i].Name < voices[j].Name
	})
	return voices
}

// speechEnsureAuthorized checks a permission and prompts when it has not
// been decided yet. It blocks until the user answers the prompt.
func speechEnsureAuthorized(kind int, status func() int, request func(), what string) error {
	current := status()
	if current == C.WailsSpeechAuthNotDetermined {
		// Drain a stale event, then prompt.
		select {
		case <-speechAuthorizationEvents:
		default:
		}
		request()
		for {
			event := <-speechAuthorizationEvents
			if event.kind == kind {
				current = event.status
				break
			}
		}
	}
	switch current {
	case C.WailsSpeechAuthAuthorized:
		return nil
	case C.WailsSpeechAuthRestricted:
		return fmt.Errorf("%w: %s is restricted by policy", ErrSpeechRecognitionDenied, what)
	default:
		return fmt.Errorf("%w: %s access was denied; enable it in System Settings > Privacy & Security", ErrSpeechRecognitionDenied, what)
	}
}

func speechRecognitionStart(id uint64, locale string) error {
	if !bool(C.speechRecognitionUsageDescriptionPresent()) || !bool(C.speechMicrophoneUsageDescriptionPresent()) {
		return ErrSpeechRecognitionUsageDescription
	}
	if err := speechEnsureAuthorized(C.WailsSpeechAuthKindRecognition,
		func() int { return int(C.speechRecognitionAuthorizationStatus()) },
		func() { C.speechRequestRecognitionAuthorization() },
		"speech recognition"); err != nil {
		return err
	}
	if err := speechEnsureAuthorized(C.WailsSpeechAuthKindMicrophone,
		func() int { return int(C.speechMicrophoneAuthorizationStatus()) },
		func() { C.speechRequestMicrophoneAuthorization() },
		"microphone"); err != nil {
		return err
	}

	var cLocale *C.char
	if locale != "" {
		cLocale = C.CString(locale)
		defer C.free(unsafe.Pointer(cLocale))
	}
	switch C.speechRecognitionStart(C.ulonglong(id), cLocale) {
	case C.WailsSpeechOK:
		return nil
	case C.WailsSpeechRecognizerUnavailable:
		if locale != "" {
			return fmt.Errorf("%w for locale %q", ErrSpeechRecognitionUnavailable, locale)
		}
		return fmt.Errorf("%w: no recogniser for the system locale", ErrSpeechRecognitionUnavailable)
	case C.WailsSpeechAudioInputUnavailable:
		return fmt.Errorf("%w: no audio input device", ErrSpeechRecognitionUnavailable)
	case C.WailsSpeechAudioEngineFailed:
		return fmt.Errorf("%w: the audio engine could not start", ErrSpeechRecognitionUnavailable)
	default:
		return ErrSpeechRecognitionNotSupported
	}
}

func speechRecognitionStop(id uint64) {
	C.speechRecognitionStop(C.ulonglong(id))
}

func speechRecognitionCancel(id uint64) {
	C.speechRecognitionCancel(C.ulonglong(id))
}
