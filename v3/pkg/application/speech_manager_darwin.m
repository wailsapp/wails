//go:build darwin && !ios && !server

#import "speech_manager_darwin.h"
#import <AVFoundation/AVFoundation.h>
#import <objc/runtime.h>
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
#import <Speech/Speech.h>
#endif

extern void speechUtteranceFinishedCallback(unsigned long long, bool);
extern void speechAuthorizationCallback(int, int);
extern void speechRecognitionResultCallback(unsigned long long, const char *, bool);
extern void speechRecognitionErrorCallback(unsigned long long, const char *);

static const void *kWailsUtteranceIDKey = &kWailsUtteranceIDKey;

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400

// WailsSpeechController owns one AVSpeechSynthesizer and its own queue of
// pending utterances. Only the current utterance is handed to the
// synthesizer so a queued utterance can be dropped individually and a stop
// of the current one does not cancel the rest of the queue. All state is
// touched on the main queue.
API_AVAILABLE(macos(10.14))
@interface WailsSpeechController : NSObject <AVSpeechSynthesizerDelegate>
@property (retain) AVSpeechSynthesizer *synthesizer;
@property (retain) NSMutableArray<AVSpeechUtterance *> *pending;
@property (retain) AVSpeechUtterance *current;
@end

@implementation WailsSpeechController

- (instancetype)init {
	self = [super init];
	if (self) {
		_synthesizer = [[AVSpeechSynthesizer alloc] init];
		_synthesizer.delegate = self;
		_pending = [[NSMutableArray alloc] init];
	}
	return self;
}

- (void)dealloc {
	_synthesizer.delegate = nil;
	[_synthesizer release];
	[_pending release];
	[_current release];
	[super dealloc];
}

+ (unsigned long long)idOfUtterance:(AVSpeechUtterance *)utterance {
	NSNumber *value = objc_getAssociatedObject(utterance, kWailsUtteranceIDKey);
	return value ? value.unsignedLongLongValue : 0;
}

- (void)enqueue:(AVSpeechUtterance *)utterance {
	[self.pending addObject:utterance];
	[self pump];
}

- (void)pump {
	if (self.current != nil || self.pending.count == 0) {
		return;
	}
	AVSpeechUtterance *next = [[self.pending objectAtIndex:0] retain];
	[self.pending removeObjectAtIndex:0];
	self.current = next;
	[next release];
	[self.synthesizer speakUtterance:next];
}

- (void)stopUtterance:(unsigned long long)utteranceID {
	if (self.current != nil && [WailsSpeechController idOfUtterance:self.current] == utteranceID) {
		[self.synthesizer stopSpeakingAtBoundary:AVSpeechBoundaryImmediate];
		return;
	}
	for (NSUInteger i = 0; i < self.pending.count; i++) {
		AVSpeechUtterance *queued = [self.pending objectAtIndex:i];
		if ([WailsSpeechController idOfUtterance:queued] == utteranceID) {
			[self.pending removeObjectAtIndex:i];
			speechUtteranceFinishedCallback(utteranceID, true);
			return;
		}
	}
}

- (void)stopAll {
	NSArray<AVSpeechUtterance *> *dropped = [[self.pending copy] autorelease];
	[self.pending removeAllObjects];
	for (AVSpeechUtterance *queued in dropped) {
		speechUtteranceFinishedCallback([WailsSpeechController idOfUtterance:queued], true);
	}
	if (self.current != nil) {
		[self.synthesizer stopSpeakingAtBoundary:AVSpeechBoundaryImmediate];
	}
}

- (void)utteranceEnded:(AVSpeechUtterance *)utterance cancelled:(BOOL)cancelled {
	dispatch_async(dispatch_get_main_queue(), ^{
		unsigned long long utteranceID = [WailsSpeechController idOfUtterance:utterance];
		if (self.current == utterance) {
			self.current = nil;
		}
		if (utteranceID != 0) {
			speechUtteranceFinishedCallback(utteranceID, (bool)cancelled);
		}
		[self pump];
	});
}

- (void)speechSynthesizer:(AVSpeechSynthesizer *)synthesizer didFinishSpeechUtterance:(AVSpeechUtterance *)utterance {
	[self utteranceEnded:utterance cancelled:NO];
}

- (void)speechSynthesizer:(AVSpeechSynthesizer *)synthesizer didCancelSpeechUtterance:(AVSpeechUtterance *)utterance {
	[self utteranceEnded:utterance cancelled:YES];
}

@end

static WailsSpeechController *wailsSpeechController(void) API_AVAILABLE(macos(10.14)) {
	static WailsSpeechController *controller = nil;
	static dispatch_once_t once;
	dispatch_once(&once, ^{
		controller = [[WailsSpeechController alloc] init];
	});
	return controller;
}

#endif // MAC_OS_X_VERSION_MAX_ALLOWED >= 101400

bool speechSynthesisAvailable(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
	if (@available(macOS 10.14, *)) {
		return true;
	}
#endif
	return false;
}

void speechSpeak(unsigned long long id, const char *text, const char *voice, double rate, double volume) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
	if (@available(macOS 10.14, *)) {
		NSString *string = text ? [NSString stringWithUTF8String:text] : @"";
		NSString *voiceID = (voice && voice[0] != '\0') ? [NSString stringWithUTF8String:voice] : nil;
		dispatch_async(dispatch_get_main_queue(), ^{
			AVSpeechUtterance *utterance = [AVSpeechUtterance speechUtteranceWithString:string];
			if (voiceID != nil) {
				AVSpeechSynthesisVoice *selected = [AVSpeechSynthesisVoice voiceWithIdentifier:voiceID];
				if (selected == nil) {
					selected = [AVSpeechSynthesisVoice voiceWithLanguage:voiceID];
				}
				if (selected != nil) {
					utterance.voice = selected;
				}
			}
			utterance.rate = (float)rate;
			utterance.volume = (float)volume;
			objc_setAssociatedObject(utterance, kWailsUtteranceIDKey, @(id), OBJC_ASSOCIATION_RETAIN_NONATOMIC);
			[wailsSpeechController() enqueue:utterance];
		});
		return;
	}
#endif
	speechUtteranceFinishedCallback(id, true);
}

void speechStopUtterance(unsigned long long id) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
	if (@available(macOS 10.14, *)) {
		dispatch_async(dispatch_get_main_queue(), ^{
			[wailsSpeechController() stopUtterance:id];
		});
	}
#endif
}

void speechStopAllUtterances(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
	if (@available(macOS 10.14, *)) {
		dispatch_async(dispatch_get_main_queue(), ^{
			[wailsSpeechController() stopAll];
		});
	}
#endif
}

char *speechVoicesJSON(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
	if (@available(macOS 10.14, *)) {
		NSMutableArray *voices = [NSMutableArray array];
		for (AVSpeechSynthesisVoice *voice in [AVSpeechSynthesisVoice speechVoices]) {
			[voices addObject:@{
				@"id": voice.identifier ?: @"",
				@"name": voice.name ?: @"",
				@"language": voice.language ?: @"",
			}];
		}
		NSData *data = [NSJSONSerialization dataWithJSONObject:voices options:0 error:nil];
		if (data == nil) {
			return NULL;
		}
		NSString *json = [[[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] autorelease];
		return strdup([json UTF8String]);
	}
#endif
	return NULL;
}

// Recognition ------------------------------------------------------------

static bool speechInfoPlistHasKey(NSString *key) {
	id value = [[NSBundle mainBundle] objectForInfoDictionaryKey:key];
	return [value isKindOfClass:[NSString class]] && [(NSString *)value length] > 0;
}

bool speechRecognitionUsageDescriptionPresent(void) {
	return speechInfoPlistHasKey(@"NSSpeechRecognitionUsageDescription");
}

bool speechMicrophoneUsageDescriptionPresent(void) {
	return speechInfoPlistHasKey(@"NSMicrophoneUsageDescription");
}

bool speechRecognitionAvailable(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
	if (@available(macOS 10.15, *)) {
		return true;
	}
#endif
	return false;
}

int speechRecognitionAuthorizationStatus(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
	if (@available(macOS 10.15, *)) {
		return (int)[SFSpeechRecognizer authorizationStatus];
	}
#endif
	return WailsSpeechAuthRestricted;
}

int speechMicrophoneAuthorizationStatus(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
	if (@available(macOS 10.14, *)) {
		return (int)[AVCaptureDevice authorizationStatusForMediaType:AVMediaTypeAudio];
	}
#endif
	return WailsSpeechAuthAuthorized;
}

void speechRequestRecognitionAuthorization(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
	if (@available(macOS 10.15, *)) {
		[SFSpeechRecognizer requestAuthorization:^(SFSpeechRecognizerAuthorizationStatus status) {
			speechAuthorizationCallback(WailsSpeechAuthKindRecognition, (int)status);
		}];
		return;
	}
#endif
	speechAuthorizationCallback(WailsSpeechAuthKindRecognition, WailsSpeechAuthRestricted);
}

void speechRequestMicrophoneAuthorization(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
	if (@available(macOS 10.14, *)) {
		[AVCaptureDevice requestAccessForMediaType:AVMediaTypeAudio completionHandler:^(BOOL granted) {
			speechAuthorizationCallback(WailsSpeechAuthKindMicrophone,
				granted ? WailsSpeechAuthAuthorized : WailsSpeechAuthDenied);
		}];
		return;
	}
#endif
	speechAuthorizationCallback(WailsSpeechAuthKindMicrophone, WailsSpeechAuthAuthorized);
}

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500

API_AVAILABLE(macos(10.15))
@interface WailsSpeechRecognition : NSObject
@property unsigned long long sessionID;
@property (retain) SFSpeechRecognizer *recognizer;
@property (retain) SFSpeechAudioBufferRecognitionRequest *request;
@property (retain) SFSpeechRecognitionTask *task;
@property (retain) AVAudioEngine *engine;
@property BOOL tapInstalled;
@property BOOL finished;
@end

@implementation WailsSpeechRecognition

- (void)dealloc {
	[_recognizer release];
	[_request release];
	[_task release];
	[_engine release];
	[super dealloc];
}

- (void)stopAudio {
	if (self.engine != nil) {
		if (self.tapInstalled) {
			[self.engine.inputNode removeTapOnBus:0];
			self.tapInstalled = NO;
		}
		[self.engine stop];
	}
	[self.request endAudio];
}

- (void)teardown {
	[self stopAudio];
	self.finished = YES;
	self.task = nil;
	self.request = nil;
	self.engine = nil;
}

@end

static NSMutableDictionary<NSNumber *, WailsSpeechRecognition *> *wailsRecognitions(void) API_AVAILABLE(macos(10.15)) {
	static NSMutableDictionary *sessions = nil;
	static dispatch_once_t once;
	dispatch_once(&once, ^{
		sessions = [[NSMutableDictionary alloc] init];
	});
	return sessions;
}

static int speechRecognitionStartOnMain(unsigned long long id, NSString *localeID) API_AVAILABLE(macos(10.15)) {
	SFSpeechRecognizer *recognizer = nil;
	if (localeID.length > 0) {
		recognizer = [[[SFSpeechRecognizer alloc] initWithLocale:[NSLocale localeWithLocaleIdentifier:localeID]] autorelease];
	} else {
		recognizer = [[[SFSpeechRecognizer alloc] init] autorelease];
	}
	if (recognizer == nil || !recognizer.available) {
		return WailsSpeechRecognizerUnavailable;
	}

	AVAudioEngine *engine = [[[AVAudioEngine alloc] init] autorelease];
	AVAudioInputNode *input = engine.inputNode;
	AVAudioFormat *format = [input outputFormatForBus:0];
	if (input == nil || format == nil || format.channelCount == 0 || format.sampleRate == 0) {
		return WailsSpeechAudioInputUnavailable;
	}

	SFSpeechAudioBufferRecognitionRequest *request = [[[SFSpeechAudioBufferRecognitionRequest alloc] init] autorelease];
	request.shouldReportPartialResults = YES;

	WailsSpeechRecognition *session = [[[WailsSpeechRecognition alloc] init] autorelease];
	session.sessionID = id;
	session.recognizer = recognizer;
	session.request = request;
	session.engine = engine;

	[input installTapOnBus:0 bufferSize:1024 format:format block:^(AVAudioPCMBuffer *buffer, AVAudioTime *when) {
		[request appendAudioPCMBuffer:buffer];
	}];
	session.tapInstalled = YES;
	[engine prepare];
	NSError *error = nil;
	if (![engine startAndReturnError:&error]) {
		[session teardown];
		return WailsSpeechAudioEngineFailed;
	}

	[wailsRecognitions() setObject:session forKey:@(id)];

	session.task = [recognizer recognitionTaskWithRequest:request resultHandler:^(SFSpeechRecognitionResult *result, NSError *taskError) {
		NSString *text = result ? result.bestTranscription.formattedString : nil;
		BOOL isFinal = result ? result.isFinal : NO;
		NSString *message = taskError ? [NSString stringWithFormat:@"%@ (code=%ld)", taskError.localizedDescription, (long)taskError.code] : nil;
		dispatch_async(dispatch_get_main_queue(), ^{
			WailsSpeechRecognition *live = [wailsRecognitions() objectForKey:@(id)];
			if (live == nil || live.finished) {
				return;
			}
			if (text != nil) {
				speechRecognitionResultCallback(id, [text UTF8String], (bool)isFinal);
			}
			if (message != nil) {
				speechRecognitionErrorCallback(id, [message UTF8String]);
			}
			if (message != nil || isFinal) {
				[live teardown];
				[wailsRecognitions() removeObjectForKey:@(id)];
			}
		});
	}];
	return WailsSpeechOK;
}

#endif // MAC_OS_X_VERSION_MAX_ALLOWED >= 101500

int speechRecognitionStart(unsigned long long id, const char *locale) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
	if (@available(macOS 10.15, *)) {
		NSString *localeID = locale ? [NSString stringWithUTF8String:locale] : nil;
		__block int result = WailsSpeechUnsupported;
		if ([NSThread isMainThread]) {
			result = speechRecognitionStartOnMain(id, localeID);
		} else {
			dispatch_sync(dispatch_get_main_queue(), ^{
				result = speechRecognitionStartOnMain(id, localeID);
			});
		}
		return result;
	}
#endif
	return WailsSpeechUnsupported;
}

void speechRecognitionStop(unsigned long long id) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
	if (@available(macOS 10.15, *)) {
		dispatch_async(dispatch_get_main_queue(), ^{
			WailsSpeechRecognition *session = [wailsRecognitions() objectForKey:@(id)];
			// endAudio lets the recogniser deliver its final result, which
			// tears the session down in the result handler.
			[session stopAudio];
		});
	}
#endif
}

void speechRecognitionCancel(unsigned long long id) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
	if (@available(macOS 10.15, *)) {
		dispatch_async(dispatch_get_main_queue(), ^{
			WailsSpeechRecognition *session = [wailsRecognitions() objectForKey:@(id)];
			if (session == nil) {
				return;
			}
			[session.task cancel];
			[session teardown];
			[wailsRecognitions() removeObjectForKey:@(id)];
		});
	}
#endif
}
