//go:build linux

package ui

import "os"

// platformExit terminates the process.  On Linux the standard os.Exit is fine
// because there are no Metal/CoreGraphics global destructors that could assert.
func platformExit() {
	os.Exit(0)
}
