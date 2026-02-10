package fileutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExpandPath resolves ~ to the user's home directory.
func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("expanding path %q: %w", path, err)
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

// IsHidden returns true if the file name starts with a dot.
func IsHidden(name string) bool {
	return strings.HasPrefix(filepath.Base(name), ".")
}

// partialExtensions are file extensions indicating incomplete downloads.
var partialExtensions = map[string]bool{
	".crdownload": true,
	".part":       true,
	".tmp":        true,
	".download":   true,
	".partial":    true,
}

// ShouldIgnore returns true if the file should not be organized.
func ShouldIgnore(path string) bool {
	name := filepath.Base(path)
	// Inline hidden check — name already has no path separators.
	if len(name) > 0 && name[0] == '.' {
		return true
	}
	ext := strings.ToLower(filepath.Ext(name))
	return partialExtensions[ext]
}

// ResolveConflict returns a non-conflicting file path by appending a timestamp
// suffix if the destination already exists.
func ResolveConflict(dest string) string {
	if _, err := os.Lstat(dest); os.IsNotExist(err) {
		return dest
	}

	dir := filepath.Dir(dest)
	ext := filepath.Ext(dest)
	base := strings.TrimSuffix(filepath.Base(dest), ext)
	stamp := time.Now().Format("20060102_150405")
	newName := fmt.Sprintf("%s_%s%s", base, stamp, ext)
	newPath := filepath.Join(dir, newName)

	// In the unlikely event the timestamped name also exists, add a counter.
	if _, err := os.Lstat(newPath); err == nil {
		for i := 2; ; i++ {
			newName = fmt.Sprintf("%s_%s_%d%s", base, stamp, i, ext)
			newPath = filepath.Join(dir, newName)
			if _, err := os.Lstat(newPath); os.IsNotExist(err) {
				break
			}
		}
	}

	return newPath
}

// EnsureDir creates the directory (and parents) if it doesn't exist.
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// MoveFile moves src to dst. It tries os.Rename first (atomic, same filesystem),
// and falls back to copy+delete for cross-device moves.
func MoveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// Fallback: copy then remove.
	if err := copyFile(src, dst); err != nil {
		return fmt.Errorf("copy fallback: %w", err)
	}
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("removing source after copy: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
