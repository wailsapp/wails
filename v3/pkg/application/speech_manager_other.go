//go:build !darwin || ios || server

package application

func speechSynthesisSupported() bool   { return false }
func speechRecognitionSupported() bool { return false }

func speechSpeak(uint64, string, SpeechOptions) error { return ErrSpeechNotSupported }
func speechStop(uint64)                               {}
func speechStopAll()                                  {}
func speechVoices() []Voice                           { return nil }

func speechRecognitionStart(uint64, string) error { return ErrSpeechRecognitionNotSupported }
func speechRecognitionStop(uint64)                {}
func speechRecognitionCancel(uint64)              {}
