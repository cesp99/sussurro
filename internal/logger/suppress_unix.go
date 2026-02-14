//go:build !windows

package logger

import (
	"os"
	"syscall"
)

// SuppressStderr redirects stderr to /dev/null and returns a cleanup function to restore it
func SuppressStderr() func() {
	// Save current stderr
	originalStderr, err := syscall.Dup(int(os.Stderr.Fd()))
	if err != nil {
		return func() {}
	}

	// Open /dev/null
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		syscall.Close(originalStderr)
		return func() {}
	}

	// Redirect stderr to /dev/null
	// This works at the OS level, affecting C libraries too
	err = syscall.Dup2(int(devNull.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		syscall.Close(originalStderr)
		devNull.Close()
		return func() {}
	}

	return func() {
		// Restore stderr
		syscall.Dup2(originalStderr, int(os.Stderr.Fd()))
		syscall.Close(originalStderr)
		devNull.Close()
	}
}
