package application

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// DefaultSpeechRate is the rate used when SpeechOptions.Rate is 0. It maps to
// AVSpeechUtteranceDefaultSpeechRate on Apple platforms.
const DefaultSpeechRate = 0.5

// speechRecognitionFinalTimeout bounds how long RecognitionSession.Stop waits
// for the recogniser to deliver its final transcript after audio ends.
const speechRecognitionFinalTimeout = 10 * time.Second

var (
	// ErrSpeechNotSupported is returned when text to speech is unavailable
	// on the current platform.
	ErrSpeechNotSupported = errors.New("speech synthesis is not supported on this platform")
	// ErrSpeechRecognitionNotSupported is returned when speech recognition
	// is unavailable on the current platform or OS version.
	ErrSpeechRecognitionNotSupported = errors.New("speech recognition is not supported on this platform")
	// ErrSpeechRecognitionDenied is returned when the user (or a policy)
	// denied speech recognition or microphone access.
	ErrSpeechRecognitionDenied = errors.New("speech recognition authorization denied")
	// ErrSpeechRecognitionUnavailable is returned when no recogniser is
	// available for the requested locale or no audio input device exists.
	ErrSpeechRecognitionUnavailable = errors.New("speech recognition is unavailable")
	// ErrSpeechRecognitionUsageDescription is returned when the app bundle's
	// Info.plist lacks the privacy usage strings required by macOS.
	ErrSpeechRecognitionUsageDescription = errors.New("Info.plist is missing NSSpeechRecognitionUsageDescription or NSMicrophoneUsageDescription")
)

// SpeechOptions configures a single Speak call.
type SpeechOptions struct {
	// Voice is the identifier of a voice returned by SpeechManager.Voices.
	// Empty selects the default voice for the system language.
	Voice string
	// Rate is the speaking rate between 0 (slowest) and 1 (fastest). 0 uses
	// DefaultSpeechRate. Values outside the range are clamped.
	Rate float64
	// Volume is the playback volume between 0 (silent) and 1 (full). 0 uses
	// full volume; set Muted for a silent utterance. Values above 1 clamp to
	// 1 and negative values clamp to 0.
	Volume float64
	// Muted speaks silently, which is useful for warming up the engine or
	// testing timing without audio.
	Muted bool
}

// normalized applies the defaults and clamping documented on SpeechOptions.
func (o SpeechOptions) normalized() SpeechOptions {
	out := o
	out.Voice = strings.TrimSpace(o.Voice)
	switch {
	case o.Rate == 0:
		out.Rate = DefaultSpeechRate
	case o.Rate < 0:
		out.Rate = 0
	case o.Rate > 1:
		out.Rate = 1
	}
	switch {
	case o.Muted:
		out.Volume = 0
	case o.Volume == 0:
		out.Volume = 1
	case o.Volume < 0:
		out.Volume = 0
	case o.Volume > 1:
		out.Volume = 1
	}
	return out
}

// Voice describes a text to speech voice installed on the system.
type Voice struct {
	// ID is the platform identifier, for example
	// "com.apple.voice.compact.en-GB.Daniel". Pass it as SpeechOptions.Voice.
	ID string
	// Name is the display name, for example "Daniel".
	Name string
	// Language is the BCP 47 tag, for example "en-GB".
	Language string
}

// Utterance is a single piece of spoken text created by SpeechManager.Speak.
type Utterance struct {
	id   uint64
	text string

	mu        sync.Mutex
	finished  bool
	cancelled bool
	handlers  []func()
	done      chan struct{}
}

// ID returns the utterance identifier.
func (u *Utterance) ID() uint64 {
	return u.id
}

// Text returns the text being spoken.
func (u *Utterance) Text() string {
	return u.text
}

// Stop stops this utterance. If it is currently being spoken it stops
// immediately; if it is still queued behind other utterances it is dropped
// without being spoken. OnFinished handlers run in both cases.
func (u *Utterance) Stop() {
	speechStop(u.id)
}

// OnFinished registers a handler called once the utterance finishes or is
// stopped. It runs on a background goroutine. When the utterance has already
// finished the handler is called immediately on the calling goroutine.
func (u *Utterance) OnFinished(handler func()) {
	if handler == nil {
		return
	}
	u.mu.Lock()
	if u.finished {
		u.mu.Unlock()
		handler()
		return
	}
	u.handlers = append(u.handlers, handler)
	u.mu.Unlock()
}

// Done returns a channel closed when the utterance finishes or is stopped.
func (u *Utterance) Done() <-chan struct{} {
	return u.done
}

// IsFinished reports whether the utterance has finished or been stopped.
func (u *Utterance) IsFinished() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.finished
}

// WasStopped reports whether the utterance ended because of Stop or StopAll
// rather than by speaking to completion.
func (u *Utterance) WasStopped() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.cancelled
}

// finish marks the utterance complete and runs the handlers once.
func (u *Utterance) finish(cancelled bool) {
	u.mu.Lock()
	if u.finished {
		u.mu.Unlock()
		return
	}
	u.finished = true
	u.cancelled = cancelled
	handlers := u.handlers
	u.handlers = nil
	close(u.done)
	u.mu.Unlock()
	for _, handler := range handlers {
		handler()
	}
}

// RecognitionOptions configures SpeechManager.Recognize.
type RecognitionOptions struct {
	// Locale is the BCP 47 tag of the language to recognise, for example
	// "en-US". Empty uses the system locale.
	Locale string
	// OnPartial receives the transcript so far each time the recogniser
	// updates it. It runs on a background goroutine.
	OnPartial func(text string)
}

// RecognitionSession is a live microphone transcription started by
// SpeechManager.Recognize. Call Stop to end capture and read the transcript.
type RecognitionSession struct {
	id        uint64
	onPartial func(string)

	mu      sync.Mutex
	text    string
	err     error
	final   chan struct{}
	ended   bool
	stopped bool
}

// ID returns the session identifier.
func (r *RecognitionSession) ID() uint64 {
	return r.id
}

// Text returns the transcript recognised so far.
func (r *RecognitionSession) Text() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.text
}

// Stop ends microphone capture, waits (up to ten seconds) for the final
// transcript and returns it. An error is returned only when the session
// failed before producing any text; a session that ended without detecting
// speech returns an empty transcript and no error. Calling Stop more than
// once returns the same result.
func (r *RecognitionSession) Stop() (string, error) {
	r.mu.Lock()
	alreadyStopped := r.stopped
	r.stopped = true
	ended := r.ended
	r.mu.Unlock()

	if !alreadyStopped && !ended {
		speechRecognitionStop(r.id)
	}
	select {
	case <-r.final:
	case <-time.After(speechRecognitionFinalTimeout):
		speechRecognitionCancel(r.id)
		r.end("", nil, true)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil && r.text == "" {
		return "", r.err
	}
	return r.text, nil
}

// Cancel abandons the session immediately without waiting for a transcript.
func (r *RecognitionSession) Cancel() {
	r.mu.Lock()
	r.stopped = true
	ended := r.ended
	r.mu.Unlock()
	if !ended {
		speechRecognitionCancel(r.id)
	}
	r.end("", nil, true)
}

// update records a new transcript. When final is true the session ends.
func (r *RecognitionSession) update(text string, final bool) {
	r.mu.Lock()
	if r.ended {
		r.mu.Unlock()
		return
	}
	changed := text != r.text
	if text != "" || final {
		r.text = text
	}
	handler := r.onPartial
	r.mu.Unlock()
	if changed && handler != nil {
		handler(text)
	}
	if final {
		r.end(text, nil, true)
	}
}

// end closes the session with the given result. No speech errors are not
// treated as failures so Stop can return the (empty) transcript cleanly.
func (r *RecognitionSession) end(text string, err error, keepText bool) {
	r.mu.Lock()
	if r.ended {
		r.mu.Unlock()
		return
	}
	r.ended = true
	if !keepText || text != "" {
		r.text = text
	}
	if err != nil && !isNoSpeechError(err) {
		r.err = err
	}
	close(r.final)
	r.mu.Unlock()
	speechRecognitionForget(r.id)
}

// isNoSpeechError reports whether err is the recogniser's "nothing was said"
// outcome rather than a failure.
func isNoSpeechError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "no speech detected") || strings.Contains(msg, "code=1110") || strings.Contains(msg, "code=203")
}

// speech registries shared by the platform implementations.
var (
	speechNextID        uint64
	speechUtterances    sync.Map // uint64 -> *Utterance
	speechRecognitions  sync.Map // uint64 -> *RecognitionSession
	speechRecognizeLock sync.Mutex
)

func nextSpeechID() uint64 {
	return atomic.AddUint64(&speechNextID, 1)
}

func speechUtteranceByID(id uint64) *Utterance {
	value, ok := speechUtterances.Load(id)
	if !ok {
		return nil
	}
	return value.(*Utterance)
}

func speechRecognitionByID(id uint64) *RecognitionSession {
	value, ok := speechRecognitions.Load(id)
	if !ok {
		return nil
	}
	return value.(*RecognitionSession)
}

func speechRecognitionForget(id uint64) {
	speechRecognitions.Delete(id)
}

// SpeechManager speaks text and recognises speech using platform services.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type SpeechManager struct {
	app *App
}

func newSpeechManager(app *App) *SpeechManager {
	return &SpeechManager{app: app}
}

// Speak queues text for speech and returns immediately. Utterances are
// spoken one after another in the order they were queued.
//
// macOS 10.14 and later: AVSpeechSynthesizer. Other platforms return
// ErrSpeechNotSupported.
func (s *SpeechManager) Speak(text string, opts SpeechOptions) (*Utterance, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("speech text is empty")
	}
	if !speechSynthesisSupported() {
		return nil, ErrSpeechNotSupported
	}
	utterance := &Utterance{
		id:   nextSpeechID(),
		text: text,
		done: make(chan struct{}),
	}
	speechUtterances.Store(utterance.id, utterance)
	if err := speechSpeak(utterance.id, text, opts.normalized()); err != nil {
		speechUtterances.Delete(utterance.id)
		return nil, err
	}
	return utterance, nil
}

// Voices returns the installed text to speech voices sorted by language and
// name. It returns nil on platforms without speech synthesis.
func (s *SpeechManager) Voices() []Voice {
	return speechVoices()
}

// StopAll stops the current utterance and drops any queued ones. Their
// OnFinished handlers run.
func (s *SpeechManager) StopAll() {
	speechStopAll()
}

// IsSupported reports whether Speak works on this platform.
func (s *SpeechManager) IsSupported() bool {
	return speechSynthesisSupported()
}

// IsRecognitionSupported reports whether Recognize works on this platform
// and OS version. It does not check authorization.
func (s *SpeechManager) IsRecognitionSupported() bool {
	return speechRecognitionSupported()
}

// Recognize starts transcribing audio from the default microphone and
// returns a session. Partial transcripts arrive through
// RecognitionOptions.OnPartial; call RecognitionSession.Stop to end capture
// and get the final text. Only one session runs at a time; starting another
// while one is active returns an error.
//
// macOS 10.15 and later: SFSpeechRecognizer with AVAudioEngine capture. The
// app must be a bundle whose Info.plist declares
// NSSpeechRecognitionUsageDescription and NSMicrophoneUsageDescription, or
// macOS refuses access (ErrSpeechRecognitionUsageDescription). The first
// call prompts the user for both permissions and blocks until they answer,
// so call Recognize from a goroutine rather than the main thread; a refusal
// returns ErrSpeechRecognitionDenied and can be changed later in System
// Settings > Privacy & Security. On-device recognition is used when the
// locale supports it, otherwise audio is sent to Apple's servers. Other
// platforms return ErrSpeechRecognitionNotSupported.
func (s *SpeechManager) Recognize(opts RecognitionOptions) (*RecognitionSession, error) {
	if !speechRecognitionSupported() {
		return nil, ErrSpeechRecognitionNotSupported
	}
	speechRecognizeLock.Lock()
	defer speechRecognizeLock.Unlock()

	active := false
	speechRecognitions.Range(func(_, _ any) bool {
		active = true
		return false
	})
	if active {
		return nil, errors.New("a speech recognition session is already running")
	}

	session := &RecognitionSession{
		id:        nextSpeechID(),
		onPartial: opts.OnPartial,
		final:     make(chan struct{}),
	}
	speechRecognitions.Store(session.id, session)
	if err := speechRecognitionStart(session.id, strings.TrimSpace(opts.Locale)); err != nil {
		speechRecognitions.Delete(session.id)
		return nil, fmt.Errorf("speech recognition: %w", err)
	}
	return session, nil
}
