package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	appLabel = "com.geraldbahati.file-organizer"
)

// PIDPath returns the path to the PID file.
func PIDPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "file-organizer")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "file-organizer.pid"), nil
}

// WritePID writes the current process PID to the PID file.
func WritePID() error {
	path, err := PIDPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strconv.Itoa(os.Getpid())), 0644)
}

// RemovePID removes the PID file.
func RemovePID() error {
	path, err := PIDPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}

// ReadPID reads the PID from the PID file. Returns 0 if not found.
func ReadPID() (int, error) {
	path, err := PIDPath()
	if err != nil {
		return 0, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("invalid PID file: %w", err)
	}
	return pid, nil
}

// IsRunning checks if a process with the given PID is running.
func IsRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix, FindProcess always succeeds. Send signal 0 to check.
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// StopDaemon sends SIGTERM to the running daemon.
func StopDaemon() error {
	pid, err := ReadPID()
	if err != nil {
		return fmt.Errorf("reading PID: %w", err)
	}
	if pid == 0 || !IsRunning(pid) {
		return fmt.Errorf("daemon is not running")
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGTERM)
}

// LaunchdPlistPath returns the path to the launchd plist file.
func LaunchdPlistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", appLabel+".plist"), nil
}

// GeneratePlist creates a launchd plist for the daemon.
func GeneratePlist(binaryPath string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>start</string>
        <string>--foreground</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <false/>
    <key>StandardOutPath</key>
    <string>/tmp/file-organizer.out.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/file-organizer.err.log</string>
</dict>
</plist>
`, appLabel, binaryPath)
}
