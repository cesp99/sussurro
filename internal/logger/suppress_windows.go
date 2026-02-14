//go:build windows

package logger

// SuppressStderr for Windows
// Currently a no-op or simple Go-level redirection as full C-level redirection
// on Windows requires more complex handle manipulation (SetStdHandle).
func SuppressStderr() func() {
	// For now, we just return a no-op on Windows to ensure it compiles.
	// Implementing robust stderr suppression on Windows for C++ libs requires
	// interacting with the CRT which is non-trivial from Go.
	return func() {}
}
