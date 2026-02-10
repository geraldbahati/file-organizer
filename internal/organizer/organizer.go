package organizer

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/geraldbahati/file-organizer/internal/fileutil"
)

// Organizer processes files and moves them to the appropriate destination.
type Organizer struct {
	index  *RuleIndex
	logger *slog.Logger
}

// New creates an Organizer with the given rule index and logger.
func New(index *RuleIndex, logger *slog.Logger) *Organizer {
	return &Organizer{
		index:  index,
		logger: logger,
	}
}

// ProcessFile checks the file against rules and moves it if a match is found.
// Returns true if the file was moved.
func (o *Organizer) ProcessFile(path string) (bool, error) {
	if fileutil.ShouldIgnore(path) {
		o.logger.Debug("ignoring file", "path", path)
		return false, nil
	}

	ext := filepath.Ext(path)
	if ext == "" {
		o.logger.Debug("no extension, skipping", "path", path)
		return false, nil
	}

	destDir := o.index.Match(ext)
	if destDir == "" {
		o.logger.Debug("no rule matched", "path", path, "ext", ext)
		return false, nil
	}

	expandedDir, err := fileutil.ExpandPath(destDir)
	if err != nil {
		return false, fmt.Errorf("expanding destination %q: %w", destDir, err)
	}

	if err := fileutil.EnsureDir(expandedDir); err != nil {
		return false, fmt.Errorf("creating destination %q: %w", expandedDir, err)
	}

	destPath := filepath.Join(expandedDir, filepath.Base(path))
	destPath = fileutil.ResolveConflict(destPath)

	o.logger.Info("moving file", "src", path, "dst", destPath)

	if err := fileutil.MoveFile(path, destPath); err != nil {
		return false, fmt.Errorf("moving %q to %q: %w", path, destPath, err)
	}

	return true, nil
}
