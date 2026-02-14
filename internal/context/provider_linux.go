//go:build linux

package context

// NewProvider creates the appropriate context provider for the current platform
func NewProvider() Provider {
	return NewLinuxProvider()
}
