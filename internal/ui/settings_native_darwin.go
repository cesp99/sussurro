//go:build darwin

package ui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

static void show_window(void *win) {
    NSWindow *w = (__bridge NSWindow *)win;
    [w makeKeyAndOrderFront:nil];
}
static void hide_window(void *win) {
    NSWindow *w = (__bridge NSWindow *)win;
    [w orderOut:nil];
}
*/
import "C"
import "unsafe"

func showWebviewWindow(win unsafe.Pointer) {
	C.show_window(win)
}

func hideWebviewWindow(win unsafe.Pointer) {
	C.hide_window(win)
}

// interceptSettingsClose is handled by the NSWindowDelegate on macOS.
func interceptSettingsClose(_ unsafe.Pointer) {}
