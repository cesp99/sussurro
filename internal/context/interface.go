package context

import (
	"fmt"
	"time"
)

// ContextInfo holds information about the current user context
type ContextInfo struct {
	AppName     string
	WindowTitle string
	Timestamp   time.Time
}

// String returns a formatted string of the context
func (c ContextInfo) String() string {
	return fmt.Sprintf("[%s] App: %s, Window: %s", c.Timestamp.Format(time.RFC3339), c.AppName, c.WindowTitle)
}

// Provider defines the interface for context detection
type Provider interface {
	GetContext() (*ContextInfo, error)
	Close() error
}
