//go:build linux

package context

import (
	"os/exec"
	"strings"
	"time"
)

// LinuxProvider implements context detection for Linux systems using X11
type LinuxProvider struct{}

// NewLinuxProvider creates a new Linux context provider
func NewLinuxProvider() *LinuxProvider {
	return &LinuxProvider{}
}

// GetContext retrieves the current application name and window title
func (p *LinuxProvider) GetContext() (*ContextInfo, error) {
	info := &ContextInfo{
		Timestamp: time.Now(),
	}

	// Try to get active window information using xdotool and xprop
	// This works for X11 sessions
	windowID, err := exec.Command("xdotool", "getactivewindow").Output()
	if err != nil {
		// If xdotool fails, return empty context
		info.AppName = "unknown"
		info.WindowTitle = "unknown"
		return info, nil
	}

	// Get window class (application name)
	classCmd := exec.Command("xprop", "-id", strings.TrimSpace(string(windowID)), "WM_CLASS")
	classOutput, err := classCmd.Output()
	if err == nil {
		// WM_CLASS returns something like: WM_CLASS(STRING) = "navigator", "Firefox"
		// We want the second value (the application name)
		classStr := string(classOutput)
		if idx := strings.LastIndex(classStr, "\""); idx > 0 {
			if idx2 := strings.LastIndex(classStr[:idx], "\""); idx2 >= 0 {
				info.AppName = classStr[idx2+1 : idx]
			}
		}
	}

	// Get window title
	titleCmd := exec.Command("xdotool", "getactivewindow", "getwindowname")
	titleOutput, err := titleCmd.Output()
	if err == nil {
		info.WindowTitle = strings.TrimSpace(string(titleOutput))
	}

	// If we couldn't get the app name or window title, set defaults
	if info.AppName == "" {
		info.AppName = "unknown"
	}
	if info.WindowTitle == "" {
		info.WindowTitle = "unknown"
	}

	return info, nil
}

// Close cleans up any resources used by the provider
func (p *LinuxProvider) Close() error {
	// No resources to clean up for Linux implementation
	return nil
}
