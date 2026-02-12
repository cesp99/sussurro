package injection

import (
	"fmt"
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
)

// Injector handles text injection into the active application
type Injector struct {
	kb keybd_event.KeyBonding
}

// NewInjector creates a new text injector
func NewInjector() (*Injector, error) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return nil, fmt.Errorf("failed to create key bonding: %w", err)
	}

	// For Linux, we might need to set the device path, but for now we assume standard setup
	if runtime.GOOS == "linux" {
		// potential TODO: handle linux specific setup
	}

	return &Injector{kb: kb}, nil
}

// TypeString simulates typing the given string
// Note: keybd_event is limited in character support. 
// For complex text, we might want to use clipboard + paste.
func (i *Injector) TypeString(text string) error {
	// For macOS, we can use the clipboard approach for reliability with special chars
	// 1. Copy text to clipboard
	// 2. Simulate Cmd+V
	
	// However, the interface asks for injection. 
	// Let's implement a hybrid approach:
	// If the text is simple, maybe type it? 
	// Actually, pasting is almost always faster and more reliable for blocks of text.
	
	return i.Paste()
}

// Paste simulates the paste command (Cmd+V on Mac, Ctrl+V on others)
func (i *Injector) Paste() error {
	// Set keys for Paste command
	if runtime.GOOS == "darwin" {
		i.kb.SetKeys(keybd_event.VK_V)
		i.kb.HasSuper(true)
	} else {
		i.kb.SetKeys(keybd_event.VK_V)
		i.kb.HasCTRL(true)
	}

	// Press and Release
	err := i.kb.Launching()
	if err != nil {
		return fmt.Errorf("failed to simulate paste: %w", err)
	}
	
	// Reset modifiers
	i.kb.HasSuper(false)
	i.kb.HasCTRL(false)
	
	return nil
}

// Inject simulates typing or pasting text.
// Currently defaults to clipboard paste as it's more robust for AI output.
func (i *Injector) Inject(text string) error {
	// We assume the caller has already placed the text in the clipboard 
	// or we should handle it here. 
	// Ideally, this package should handle the clipboard write too to be self-contained.
	// But `clipboard` is in a separate package. 
	// Let's assume the pipeline handles clipboard write, and this handles the keypress.
	// OR we import internal/clipboard here.
	
	// Let's make this method do the pasting action.
	// The pipeline currently writes to clipboard.
	
	// Add a small delay to ensure clipboard is ready and window is focused
	time.Sleep(100 * time.Millisecond)
	
	return i.Paste()
}
