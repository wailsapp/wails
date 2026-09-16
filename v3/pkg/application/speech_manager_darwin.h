//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>

// Result codes shared with speech_manager_darwin.go.
enum {
	WailsSpeechOK = 0,
	WailsSpeechUnsupported = 1,
	WailsSpeechRecognizerUnavailable = 2,
	WailsSpeechAudioInputUnavailable = 3,
	WailsSpeechAudioEngineFailed = 4,
};

// Authorization status values (mirror SFSpeechRecognizerAuthorizationStatus
// and AVAuthorizationStatus): 0 not determined, 1 denied, 2 restricted,
// 3 authorized.
enum {
	WailsSpeechAuthNotDetermined = 0,
	WailsSpeechAuthDenied = 1,
	WailsSpeechAuthRestricted = 2,
	WailsSpeechAuthAuthorized = 3,
};

// Authorization kinds passed to speechAuthorizationCallback.
enum {
	WailsSpeechAuthKindRecognition = 0,
	WailsSpeechAuthKindMicrophone = 1,
};

bool speechSynthesisAvailable(void);
void speechSpeak(unsigned long long id, const char *text, const char *voice, double rate, double volume);
void speechStopUtterance(unsigned long long id);
void speechStopAllUtterances(void);
// Returns a malloc'd JSON array of {id, name, language}; caller frees.
char *speechVoicesJSON(void);

bool speechRecognitionAvailable(void);
bool speechRecognitionUsageDescriptionPresent(void);
bool speechMicrophoneUsageDescriptionPresent(void);
int speechRecognitionAuthorizationStatus(void);
int speechMicrophoneAuthorizationStatus(void);
// Both request functions report through speechAuthorizationCallback.
void speechRequestRecognitionAuthorization(void);
void speechRequestMicrophoneAuthorization(void);
// locale may be NULL for the system locale. Returns a WailsSpeech* code.
int speechRecognitionStart(unsigned long long id, const char *locale);
void speechRecognitionStop(unsigned long long id);
void speechRecognitionCancel(unsigned long long id);
