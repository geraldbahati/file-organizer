package organizer

import (
	"strings"

	"github.com/geraldbahati/file-organizer/internal/config"
)

// RuleIndex provides O(1) extension → destination lookup.
type RuleIndex struct {
	extMap map[string]string // ".pdf" → "~/Documents"
}

// NewRuleIndex builds an index from config rules.
func NewRuleIndex(rules []config.Rule) *RuleIndex {
	m := make(map[string]string, 64)
	for _, r := range rules {
		for _, ext := range r.Extensions {
			m[strings.ToLower(ext)] = r.Destination
		}
	}
	return &RuleIndex{extMap: m}
}

// Match returns the destination directory for the given file extension,
// or empty string if no rule matches.
func (ri *RuleIndex) Match(ext string) string {
	return ri.extMap[strings.ToLower(ext)]
}
