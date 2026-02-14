//go:build darwin

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>
#import <stdlib.h>

// Writes a string to the general pasteboard
void writeToClipboard(const char* text) {
    @autoreleasepool {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        [pasteboard clearContents];
        NSString *nsText = [NSString stringWithUTF8String:text];
        [pasteboard setString:nsText forType:NSPasteboardTypeString];
    }
}

// Reads a string from the general pasteboard
char* readFromClipboard() {
    @autoreleasepool {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        NSString *text = [pasteboard stringForType:NSPasteboardTypeString];
        if (text == nil) {
            return NULL;
        }
        return strdup([text UTF8String]);
    }
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Write copies the provided string to the system clipboard
func Write(text string) error {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	C.writeToClipboard(cText)
	return nil
}

// Read returns the current string content of the system clipboard
func Read() (string, error) {
	cText := C.readFromClipboard()
	if cText == nil {
		return "", errors.New("clipboard is empty or contains non-text data")
	}
	defer C.free(unsafe.Pointer(cText))

	return C.GoString(cText), nil
}
