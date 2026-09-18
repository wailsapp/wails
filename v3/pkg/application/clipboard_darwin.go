//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#import <Cocoa/Cocoa.h>
#import <stdlib.h>

bool setClipboardText(const char* text) {
	NSPasteboard *pasteBoard = [NSPasteboard generalPasteboard];
	NSError *error = nil;
	NSString *string = [NSString stringWithUTF8String:text];
	[pasteBoard clearContents];
	return [pasteBoard setString:string forType:NSPasteboardTypeString];
}

const char* getClipboardText() {
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	NSString *text = [pasteboard stringForType:NSPasteboardTypeString];
	return [text UTF8String];
}

// Rich clipboard support. Every function below runs on the main thread and
// returns malloc'd memory that the Go side frees.

// clipboardCopyBytes copies an NSData into a malloc'd buffer.
static void* clipboardCopyBytes(NSData* data, int* outLength) {
	if (data == nil || data.length == 0) {
		*outLength = 0;
		return NULL;
	}
	void* buffer = malloc(data.length);
	if (buffer == NULL) {
		*outLength = 0;
		return NULL;
	}
	memcpy(buffer, data.bytes, data.length);
	*outLength = (int)data.length;
	return buffer;
}

// clipboardCopyStrings copies an array of NSStrings into a malloc'd array
// of malloc'd C strings.
static char** clipboardCopyStrings(NSArray<NSString*>* strings, int* outCount) {
	*outCount = 0;
	if (strings == nil || strings.count == 0) {
		return NULL;
	}
	char** result = (char**)calloc(strings.count, sizeof(char*));
	if (result == NULL) {
		return NULL;
	}
	int count = 0;
	for (NSString* string in strings) {
		const char* utf8 = [string UTF8String];
		if (utf8 == NULL) {
			continue;
		}
		result[count++] = strdup(utf8);
	}
	*outCount = count;
	return result;
}

bool clipboardSetPNG(const void* bytes, int length) {
	NSData* png = [NSData dataWithBytes:bytes length:length];
	NSImage* image = [[[NSImage alloc] initWithData:png] autorelease];
	if (image == nil) {
		return false;
	}
	NSPasteboard* pasteboard = [NSPasteboard generalPasteboard];
	[pasteboard clearContents];
	[pasteboard declareTypes:@[NSPasteboardTypePNG, NSPasteboardTypeTIFF] owner:nil];
	bool ok = [pasteboard setData:png forType:NSPasteboardTypePNG];
	NSData* tiff = [image TIFFRepresentation];
	if (tiff != nil) {
		ok = [pasteboard setData:tiff forType:NSPasteboardTypeTIFF] && ok;
	}
	return ok;
}

// clipboardPNG returns the clipboard image as PNG, converting TIFF when
// that is the only bitmap on the pasteboard.
void* clipboardPNG(int* outLength) {
	NSPasteboard* pasteboard = [NSPasteboard generalPasteboard];
	NSData* png = [pasteboard dataForType:NSPasteboardTypePNG];
	if (png == nil) {
		NSData* tiff = [pasteboard dataForType:NSPasteboardTypeTIFF];
		if (tiff != nil) {
			NSBitmapImageRep* rep = [NSBitmapImageRep imageRepWithData:tiff];
			png = [rep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
		}
	}
	return clipboardCopyBytes(png, outLength);
}

bool clipboardSetFiles(const char** paths, int count) {
	NSMutableArray* urls = [NSMutableArray arrayWithCapacity:count];
	for (int i = 0; i < count; i++) {
		NSString* path = [NSString stringWithUTF8String:paths[i]];
		if (path == nil) {
			continue;
		}
		[urls addObject:[NSURL fileURLWithPath:path]];
	}
	if (urls.count == 0) {
		return false;
	}
	NSPasteboard* pasteboard = [NSPasteboard generalPasteboard];
	[pasteboard clearContents];
	return [pasteboard writeObjects:urls];
}

char** clipboardFiles(int* outCount) {
	NSPasteboard* pasteboard = [NSPasteboard generalPasteboard];
	NSArray* urls = [pasteboard readObjectsForClasses:@[[NSURL class]]
		options:@{NSPasteboardURLReadingFileURLsOnlyKey: @YES}];
	NSMutableArray* paths = [NSMutableArray arrayWithCapacity:urls.count];
	for (NSURL* url in urls) {
		if (url.path != nil) {
			[paths addObject:url.path];
		}
	}
	return clipboardCopyStrings(paths, outCount);
}

bool clipboardSetHTML(const char* html, const char* plain) {
	NSString* htmlString = [NSString stringWithUTF8String:html];
	if (htmlString == nil) {
		return false;
	}
	NSPasteboard* pasteboard = [NSPasteboard generalPasteboard];
	[pasteboard clearContents];
	NSMutableArray* types = [NSMutableArray arrayWithObject:NSPasteboardTypeHTML];
	NSString* plainString = plain != NULL && plain[0] != '\0' ? [NSString stringWithUTF8String:plain] : nil;
	if (plainString != nil) {
		[types addObject:NSPasteboardTypeString];
	}
	[pasteboard declareTypes:types owner:nil];
	bool ok = [pasteboard setString:htmlString forType:NSPasteboardTypeHTML];
	if (plainString != nil) {
		ok = [pasteboard setString:plainString forType:NSPasteboardTypeString] && ok;
	}
	return ok;
}

char* clipboardHTML() {
	NSPasteboard* pasteboard = [NSPasteboard generalPasteboard];
	NSString* html = [pasteboard stringForType:NSPasteboardTypeHTML];
	if (html == nil) {
		return NULL;
	}
	return strdup([html UTF8String]);
}

bool clipboardSetData(const char* uti, const void* bytes, int length) {
	NSString* type = [NSString stringWithUTF8String:uti];
	if (type == nil) {
		return false;
	}
	NSPasteboard* pasteboard = [NSPasteboard generalPasteboard];
	[pasteboard clearContents];
	[pasteboard declareTypes:@[type] owner:nil];
	NSData* data = length > 0 ? [NSData dataWithBytes:bytes length:length] : [NSData data];
	return [pasteboard setData:data forType:type];
}

void* clipboardData(const char* uti, int* outLength) {
	NSString* type = [NSString stringWithUTF8String:uti];
	if (type == nil) {
		*outLength = 0;
		return NULL;
	}
	NSPasteboard* pasteboard = [NSPasteboard generalPasteboard];
	return clipboardCopyBytes([pasteboard dataForType:type], outLength);
}

bool clipboardHasType(const char* uti) {
	NSString* type = [NSString stringWithUTF8String:uti];
	if (type == nil) {
		return false;
	}
	return [[[NSPasteboard generalPasteboard] types] containsObject:type];
}

char** clipboardTypes(int* outCount) {
	return clipboardCopyStrings([[NSPasteboard generalPasteboard] types], outCount);
}

void clipboardClear() {
	[[NSPasteboard generalPasteboard] clearContents];
}

long clipboardChangeCount() {
	return (long)[[NSPasteboard generalPasteboard] changeCount];
}

*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

const (
	// clipboardUTIRTF is the UTI AppKit uses for NSPasteboardTypeRTF.
	clipboardUTIRTF = "public.rtf"
)

// clipboardTakeStrings frees an array returned by clipboardCopyStrings and
// returns its contents as Go strings.
func clipboardTakeStrings(array **C.char, count C.int) []string {
	if array == nil || count == 0 {
		return nil
	}
	items := unsafe.Slice(array, int(count))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		result = append(result, C.GoString(item))
		C.free(unsafe.Pointer(item))
	}
	C.free(unsafe.Pointer(array))
	return result
}

// clipboardTakeBytes copies and frees a malloc'd buffer returned by the
// pasteboard helpers.
func clipboardTakeBytes(buffer unsafe.Pointer, length C.int) []byte {
	if buffer == nil {
		return nil
	}
	result := C.GoBytes(buffer, length)
	C.free(buffer)
	return result
}

// clipboardBytesPointer returns a pointer usable for a C const void* even
// when the slice is empty.
func clipboardBytesPointer(data []byte) unsafe.Pointer {
	if len(data) == 0 {
		return nil
	}
	return unsafe.Pointer(&data[0])
}

var clipboardLock sync.RWMutex

type macosClipboard struct{}

func (m macosClipboard) setText(text string) bool {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	cText := C.CString(text)
	success := C.setClipboardText(cText)
	C.free(unsafe.Pointer(cText))
	return bool(success)
}

func (m macosClipboard) text() (string, bool) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	clipboardText := C.getClipboardText()
	result := C.GoString(clipboardText)
	return result, true
}

func newClipboardImpl() *macosClipboard {
	return &macosClipboard{}
}

// macosClipboard also satisfies clipboardExtendedImpl. The methods run on
// the main thread (Clipboard dispatches them) and hold clipboardLock like
// the text methods so mixed callers serialise.

func (m macosClipboard) setImage(png []byte) error {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	if !bool(C.clipboardSetPNG(clipboardBytesPointer(png), C.int(len(png)))) {
		return errors.New("clipboard: the data is not a decodable PNG image")
	}
	return nil
}

func (m macosClipboard) image() ([]byte, error) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	var length C.int
	result := clipboardTakeBytes(C.clipboardPNG(&length), length)
	if result == nil {
		return nil, errors.New("clipboard: no image on the clipboard")
	}
	return result, nil
}

func (m macosClipboard) setFiles(paths []string) error {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	cPaths := make([]*C.char, len(paths))
	for i, path := range paths {
		cPaths[i] = C.CString(path)
	}
	defer func() {
		for _, cPath := range cPaths {
			C.free(unsafe.Pointer(cPath))
		}
	}()
	if !bool(C.clipboardSetFiles((**C.char)(unsafe.Pointer(&cPaths[0])), C.int(len(cPaths)))) {
		return errors.New("clipboard: unable to write the file references")
	}
	return nil
}

func (m macosClipboard) files() ([]string, error) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	var count C.int
	return clipboardTakeStrings(C.clipboardFiles(&count), count), nil
}

func (m macosClipboard) setHTML(html string, plain string) error {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	cHTML := C.CString(html)
	defer C.free(unsafe.Pointer(cHTML))
	cPlain := C.CString(plain)
	defer C.free(unsafe.Pointer(cPlain))
	if !bool(C.clipboardSetHTML(cHTML, cPlain)) {
		return errors.New("clipboard: unable to write the HTML")
	}
	return nil
}

func (m macosClipboard) html() (string, error) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	cHTML := C.clipboardHTML()
	if cHTML == nil {
		return "", errors.New("clipboard: no HTML on the clipboard")
	}
	defer C.free(unsafe.Pointer(cHTML))
	return C.GoString(cHTML), nil
}

func (m macosClipboard) setRTF(rtf []byte) error {
	return m.setData(clipboardUTIRTF, rtf)
}

func (m macosClipboard) rtf() ([]byte, error) {
	return m.data(clipboardUTIRTF)
}

func (m macosClipboard) setData(uti string, data []byte) error {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	cUTI := C.CString(uti)
	defer C.free(unsafe.Pointer(cUTI))
	if !bool(C.clipboardSetData(cUTI, clipboardBytesPointer(data), C.int(len(data)))) {
		return errors.New("clipboard: unable to write data for type " + uti)
	}
	return nil
}

func (m macosClipboard) data(uti string) ([]byte, error) {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	cUTI := C.CString(uti)
	defer C.free(unsafe.Pointer(cUTI))
	if !bool(C.clipboardHasType(cUTI)) {
		return nil, errors.New("clipboard: no data of type " + uti + " on the clipboard")
	}
	var length C.int
	result := clipboardTakeBytes(C.clipboardData(cUTI, &length), length)
	if result == nil {
		result = []byte{}
	}
	return result, nil
}

func (m macosClipboard) types() []string {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	var count C.int
	return clipboardTakeStrings(C.clipboardTypes(&count), count)
}

func (m macosClipboard) clear() {
	clipboardLock.Lock()
	defer clipboardLock.Unlock()
	C.clipboardClear()
}

func (m macosClipboard) changeCount() int {
	clipboardLock.RLock()
	defer clipboardLock.RUnlock()
	return int(C.clipboardChangeCount())
}
