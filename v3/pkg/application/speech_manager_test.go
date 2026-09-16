package application

import (
	"errors"
	"testing"
	"time"
)

func TestSpeechOptionsNormalized(t *testing.T) {
	tests := []struct {
		name string
		in   SpeechOptions
		want SpeechOptions
	}{
		{name: "zero value uses defaults", in: SpeechOptions{}, want: SpeechOptions{Rate: DefaultSpeechRate, Volume: 1}},
		{name: "explicit values kept", in: SpeechOptions{Voice: " v1 ", Rate: 0.25, Volume: 0.5}, want: SpeechOptions{Voice: "v1", Rate: 0.25, Volume: 0.5}},
		{name: "rate clamps high", in: SpeechOptions{Rate: 5}, want: SpeechOptions{Rate: 1, Volume: 1}},
		{name: "rate clamps low", in: SpeechOptions{Rate: -1}, want: SpeechOptions{Rate: 0, Volume: 1}},
		{name: "volume clamps high", in: SpeechOptions{Volume: 2}, want: SpeechOptions{Rate: DefaultSpeechRate, Volume: 1}},
		{name: "volume clamps low", in: SpeechOptions{Volume: -3}, want: SpeechOptions{Rate: DefaultSpeechRate, Volume: 0}},
		{name: "muted silences", in: SpeechOptions{Volume: 0.8, Muted: true}, want: SpeechOptions{Rate: DefaultSpeechRate, Volume: 0, Muted: true}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.in.normalized(); got != test.want {
				t.Errorf("normalized() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestUtteranceFinishRunsHandlersOnce(t *testing.T) {
	utterance := &Utterance{id: 1, text: "hi", done: make(chan struct{})}
	calls := 0
	utterance.OnFinished(func() { calls++ })
	utterance.OnFinished(nil)
	if utterance.IsFinished() {
		t.Fatal("utterance should not be finished yet")
	}
	utterance.finish(true)
	utterance.finish(false)
	if calls != 1 {
		t.Errorf("handler calls = %d, want 1", calls)
	}
	if !utterance.IsFinished() || !utterance.WasStopped() {
		t.Errorf("finished=%v stopped=%v, want both true", utterance.IsFinished(), utterance.WasStopped())
	}
	select {
	case <-utterance.Done():
	default:
		t.Error("Done channel should be closed")
	}
	// Late registration runs immediately.
	late := false
	utterance.OnFinished(func() { late = true })
	if !late {
		t.Error("OnFinished after completion should run immediately")
	}
	if utterance.ID() != 1 || utterance.Text() != "hi" {
		t.Errorf("ID/Text = %d/%q", utterance.ID(), utterance.Text())
	}
}

func TestSpeechManagerSpeakRejectsEmptyText(t *testing.T) {
	manager := newSpeechManager(nil)
	if _, err := manager.Speak("   ", SpeechOptions{}); err == nil {
		t.Error("Speak with empty text should fail")
	}
}

func TestRecognitionSessionPartialsAndStop(t *testing.T) {
	var partials []string
	session := &RecognitionSession{
		id:        7,
		onPartial: func(text string) { partials = append(partials, text) },
		final:     make(chan struct{}),
	}
	speechRecognitions.Store(session.id, session)

	session.update("hello", false)
	session.update("hello", false) // unchanged: no callback
	session.update("hello world", false)
	if session.Text() != "hello world" {
		t.Errorf("Text = %q", session.Text())
	}
	if len(partials) != 2 {
		t.Errorf("partials = %v, want two distinct updates", partials)
	}

	session.update("hello world.", true)
	text, err := session.Stop()
	if err != nil || text != "hello world." {
		t.Errorf("Stop = %q, %v", text, err)
	}
	if _, ok := speechRecognitions.Load(session.id); ok {
		t.Error("session should be forgotten after ending")
	}
	// Late events are ignored once ended.
	session.update("ignored", false)
	if session.Text() != "hello world." {
		t.Errorf("Text changed after end: %q", session.Text())
	}
	if text, err := session.Stop(); err != nil || text != "hello world." {
		t.Errorf("second Stop = %q, %v", text, err)
	}
}

func TestRecognitionSessionErrorHandling(t *testing.T) {
	failed := &RecognitionSession{id: 8, final: make(chan struct{})}
	failed.end("", errors.New("The operation couldn't be completed. (code=1101)"), true)
	if _, err := failed.Stop(); err == nil {
		t.Error("Stop should surface a failure with no transcript")
	}

	noSpeech := &RecognitionSession{id: 9, final: make(chan struct{})}
	noSpeech.end("", errors.New("No speech detected (code=1110)"), true)
	if text, err := noSpeech.Stop(); err != nil || text != "" {
		t.Errorf("no speech should be a clean empty result, got %q, %v", text, err)
	}

	withText := &RecognitionSession{id: 10, final: make(chan struct{})}
	withText.update("partial", false)
	withText.end(withText.Text(), errors.New("boom (code=1)"), true)
	if text, err := withText.Stop(); err != nil || text != "partial" {
		t.Errorf("errors after text should return the text, got %q, %v", text, err)
	}
}

func TestRecognitionSessionCancel(t *testing.T) {
	session := &RecognitionSession{id: 11, final: make(chan struct{})}
	session.update("something", false)
	done := make(chan struct{})
	go func() {
		session.Cancel()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Cancel did not return")
	}
	if text, err := session.Stop(); err != nil || text != "something" {
		t.Errorf("Stop after Cancel = %q, %v", text, err)
	}
}

func TestIsNoSpeechError(t *testing.T) {
	if isNoSpeechError(nil) {
		t.Error("nil is not a no-speech error")
	}
	if !isNoSpeechError(errors.New("No speech detected")) {
		t.Error("expected no speech detected to match")
	}
	if isNoSpeechError(errors.New("network down")) {
		t.Error("unrelated errors must not match")
	}
}
