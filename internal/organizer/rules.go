package organizer

import (
	"fmt"
	"strings"

	"github.com/geraldbahati/file-organizer/internal/config"
	"github.com/geraldbahati/file-organizer/internal/fileutil"
)

// RuleIndex provides O(1) extension → destination lookup.
// Destinations are pre-expanded absolute paths (no ~ prefix).
type RuleIndex struct {
	extMap map[string]string // ".pdf" → "/Users/x/Documents"
}

// NewRuleIndex builds an index from config rules.
// It expands ~ paths once at construction time and pre-lowers all extensions
// so Match() can do a direct map lookup with no per-call overhead.
func NewRuleIndex(rules []config.Rule) (*RuleIndex, error) {
	// Pre-compute total extension count for exact map allocation.
	n := 0
	for _, r := range rules {
		n += len(r.Extensions)
	}

	m := make(map[string]string, n)
	for _, r := range rules {
		dest, err := fileutil.ExpandPath(r.Destination)
		if err != nil {
			return nil, fmt.Errorf("expanding destination %q: %w", r.Destination, err)
		}
		for _, ext := range r.Extensions {
			m[strings.ToLower(ext)] = dest
		}
	}
	return &RuleIndex{extMap: m}, nil
}

// Match returns the destination directory for the given file extension,
// or empty string if no rule matches.
func (ri *RuleIndex) Match(ext string) string {
	return ri.extMap[strings.ToLower(ext)]
}
