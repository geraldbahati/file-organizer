package organizer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/geraldbahati/file-organizer/internal/fileutil"
)

// Organizer processes files and moves them to the appropriate destination.
type Organizer struct {
	index   *RuleIndex
	logger  *slog.Logger
	dirOnce sync.Map // destination dir → *sync.Once for EnsureDir
}

// New creates an Organizer with the given rule index and logger.
func New(index *RuleIndex, logger *slog.Logger) *Organizer {
	return &Organizer{
		index:  index,
		logger: logger,
	}
}

// ensureDir calls EnsureDir at most once per unique directory.
func (o *Organizer) ensureDir(dir string) error {
	actual, _ := o.dirOnce.LoadOrStore(dir, &sync.Once{})
	once := actual.(*sync.Once)
	var ensureErr error
	once.Do(func() {
		ensureErr = fileutil.EnsureDir(dir)
	})
	return ensureErr
}

// ProcessFile checks the file against rules and moves it if a match is found.
// Returns true if the file was moved.
func (o *Organizer) ProcessFile(path string) (bool, error) {
	// File may have already been moved by an earlier debounced event.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		o.logger.Debug("file no longer exists, skipping", "path", path)
		return false, nil
	}

	if fileutil.ShouldIgnore(path) {
		o.logger.Debug("ignoring file", "path", path)
		return false, nil
	}

	ext := filepath.Ext(path)
	if ext == "" {
		o.logger.Debug("no extension, skipping", "path", path)
		return false, nil
	}

	// Destinations are already absolute paths (expanded at RuleIndex construction).
	destDir := o.index.Match(ext)
	if destDir == "" {
		o.logger.Debug("no rule matched", "path", path, "ext", ext)
		return false, nil
	}

	if err := o.ensureDir(destDir); err != nil {
		return false, fmt.Errorf("creating destination %q: %w", destDir, err)
	}

	destPath := filepath.Join(destDir, filepath.Base(path))
	destPath = fileutil.ResolveConflict(destPath)

	o.logger.Info("moving file", "src", path, "dst", destPath)

	if err := fileutil.MoveFile(path, destPath); err != nil {
		return false, fmt.Errorf("moving %q to %q: %w", path, destPath, err)
	}

	return true, nil
}
