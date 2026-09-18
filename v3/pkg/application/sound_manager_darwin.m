//go:build darwin && !ios && !server

#import "sound_manager_darwin.h"

// WailsSoundPlayer keeps NSSound instances alive while they play and drops
// them when playback finishes. NSSound does not retain itself during
// asynchronous playback, so without this a sound released right after
// -play could be deallocated mid-stream.
@interface WailsSoundPlayer : NSObject <NSSoundDelegate>
@property (retain) NSMutableSet<NSSound *> *playing;
@end

@implementation WailsSoundPlayer

- (instancetype)init {
	self = [super init];
	if (self) {
		_playing = [[NSMutableSet alloc] init];
	}
	return self;
}

- (void)dealloc {
	[_playing release];
	[super dealloc];
}

- (int)play:(NSSound *)sound {
	if (sound == nil) {
		return WailsSoundNotFound;
	}
	sound.delegate = self;
	[self.playing addObject:sound];
	if (![sound play]) {
		[self.playing removeObject:sound];
		return WailsSoundPlayFailed;
	}
	return WailsSoundOK;
}

- (void)sound:(NSSound *)sound didFinishPlaying:(BOOL)flag {
	sound.delegate = nil;
	[self.playing removeObject:sound];
}

@end

static WailsSoundPlayer *wailsSoundPlayer(void) {
	static WailsSoundPlayer *player = nil;
	static dispatch_once_t once;
	dispatch_once(&once, ^{
		player = [[WailsSoundPlayer alloc] init];
	});
	return player;
}

void soundBeep(void) {
	NSBeep();
}

int soundPlayNamed(const char *name) {
	if (name == NULL) {
		return WailsSoundNotFound;
	}
	NSSound *sound = [NSSound soundNamed:[NSString stringWithUTF8String:name]];
	return [wailsSoundPlayer() play:sound];
}

int soundPlayFile(const char *path) {
	if (path == NULL) {
		return WailsSoundNotFound;
	}
	NSSound *sound = [[[NSSound alloc] initWithContentsOfFile:[NSString stringWithUTF8String:path] byReference:NO] autorelease];
	return [wailsSoundPlayer() play:sound];
}

int soundPlayData(const void *bytes, int length) {
	if (bytes == NULL || length <= 0) {
		return WailsSoundNotFound;
	}
	NSData *data = [NSData dataWithBytes:bytes length:(NSUInteger)length];
	NSSound *sound = [[[NSSound alloc] initWithData:data] autorelease];
	return [wailsSoundPlayer() play:sound];
}
