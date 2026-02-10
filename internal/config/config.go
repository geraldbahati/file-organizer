package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geraldbahati/file-organizer/internal/fileutil"
	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	WatchDirs  []string `yaml:"watch_dirs"`
	Rules      []Rule   `yaml:"rules"`
	DebounceMs int      `yaml:"debounce_ms"`
}

// DefaultConfigPath returns ~/.config/file-organizer/config.yaml.
func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "file-organizer", "config.yaml"), nil
}

// Load reads configuration from the given path. If the file doesn't exist,
// it creates it with defaults and returns the default config.
func Load(path string) (Config, error) {
	expanded, err := fileutil.ExpandPath(path)
	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(expanded); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if writeErr := writeDefault(expanded, cfg); writeErr != nil {
			return cfg, fmt.Errorf("creating default config: %w", writeErr)
		}
		return cfg, nil
	}

	data, err := os.ReadFile(expanded)
	if err != nil {
		return Config{}, fmt.Errorf("reading config %s: %w", expanded, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config %s: %w", expanded, err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate checks the config for common errors.
func (c *Config) Validate() error {
	if len(c.WatchDirs) == 0 {
		return fmt.Errorf("no watch directories configured")
	}
	if len(c.Rules) == 0 {
		return fmt.Errorf("no rules configured")
	}
	if c.DebounceMs < 0 {
		return fmt.Errorf("debounce_ms must be non-negative")
	}
	for _, r := range c.Rules {
		if r.Destination == "" {
			return fmt.Errorf("rule %q has no destination", r.Category)
		}
		if len(r.Extensions) == 0 {
			return fmt.Errorf("rule %q has no extensions", r.Category)
		}
	}
	return nil
}

// ExpandedWatchDirs returns watch directories with ~ expanded.
func (c *Config) ExpandedWatchDirs() ([]string, error) {
	dirs := make([]string, 0, len(c.WatchDirs))
	for _, d := range c.WatchDirs {
		expanded, err := fileutil.ExpandPath(d)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, expanded)
	}
	return dirs, nil
}

func writeDefault(path string, cfg Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	header := []byte("# file-organizer configuration\n# See https://github.com/geraldbahati/file-organizer for documentation\n\n")
	return os.WriteFile(path, append(header, data...), 0644)
}
