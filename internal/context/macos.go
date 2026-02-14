//go:build darwin

package context

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework CoreGraphics

#import <Cocoa/Cocoa.h>
#import <CoreGraphics/CoreGraphics.h>
#import <stdlib.h>

// Define a struct to hold our results
typedef struct {
    char* appName;
    char* windowTitle;
} WindowInfo;

WindowInfo getActiveWindowInfo() {
    WindowInfo info = {NULL, NULL};

    // Autorelease pool is required for Objective-C memory management
    @autoreleasepool {
        // 1. Get the frontmost (active) application using NSWorkspace
        NSRunningApplication* app = [[NSWorkspace sharedWorkspace] frontmostApplication];
        if (app == nil) {
            return info;
        }

        // Get the Application Name
        const char* name = [app.localizedName UTF8String];
        if (name) {
            info.appName = strdup(name);
        }

        pid_t pid = [app processIdentifier];

        // 2. Get a list of all visible windows on the screen
        // kCGWindowListOptionOnScreenOnly: Exclude off-screen windows
        // kCGWindowListExcludeDesktopElements: Exclude desktop icons/wallpaper
        CFArrayRef windowList = CGWindowListCopyWindowInfo(
            kCGWindowListOptionOnScreenOnly | kCGWindowListExcludeDesktopElements,
            kCGNullWindowID
        );

        if (windowList) {
            NSArray* windows = (__bridge NSArray*)windowList;

            // Iterate through windows to find the top-most window owned by our active PID
            for (NSDictionary* win in windows) {
                NSNumber* ownerPID = win[(__bridge NSString*)kCGWindowOwnerPID];

                if ([ownerPID intValue] == pid) {
                    // Check window layer. Layer 0 is the default application window layer.
                    // This filters out floating panels, status bars, etc.
                    NSNumber* layer = win[(__bridge NSString*)kCGWindowLayer];
                    if ([layer intValue] == 0) {
                        NSString* title = win[(__bridge NSString*)kCGWindowName];
                        if (title) {
                            info.windowTitle = strdup([title UTF8String]);
                        } else {
                            // If title is missing (often due to permissions), return empty string
                            info.windowTitle = strdup("");
                        }
                        // Since the list is Z-ordered, the first match is the front-most window
                        break;
                    }
                }
            }
            CFRelease(windowList);
        }
    }
    return info;
}
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

// MacOSProvider implements the Provider interface for macOS
type MacOSProvider struct{}

// NewMacOSProvider creates a new instance of MacOSProvider
func NewMacOSProvider() *MacOSProvider {
	return &MacOSProvider{}
}

// GetContext retrieves the current active window information
func (p *MacOSProvider) GetContext() (*ContextInfo, error) {
	// Call the C function
	info := C.getActiveWindowInfo()

	// Convert C strings to Go strings
	appName := C.GoString(info.appName)
	windowTitle := C.GoString(info.windowTitle)

	// Free the memory allocated by strdup in C
	if info.appName != nil {
		C.free(unsafe.Pointer(info.appName))
	}
	if info.windowTitle != nil {
		C.free(unsafe.Pointer(info.windowTitle))
	}

	if appName == "" {
		return nil, fmt.Errorf("failed to retrieve active application name")
	}

	return &ContextInfo{
		AppName:     appName,
		WindowTitle: windowTitle,
		Timestamp:   time.Now(),
	}, nil
}

// Close releases any resources (none for this implementation)
func (p *MacOSProvider) Close() error {
	return nil
}
