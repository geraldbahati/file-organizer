package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigUsesMacFriendlyDestinations(t *testing.T) {
	cfg := DefaultConfig()

	if len(cfg.WatchDirs) != 2 {
		t.Fatalf("expected 2 watch dirs, got %d", len(cfg.WatchDirs))
	}

	wantDestinations := map[string]string{
		"Archives":           "~/Downloads/Archives",
		"Disk Images":        "~/Downloads/Disk Images",
		"Installer Packages": "~/Downloads/Installers",
		"PDFs":               "~/Documents/PDFs",
		"Writing":            "~/Documents/Writing",
	}

	gotDestinations := make(map[string]string, len(cfg.Rules))
	for _, rule := range cfg.Rules {
		gotDestinations[rule.Category] = rule.Destination
	}

	for category, want := range wantDestinations {
		if got := gotDestinations[category]; got != want {
			t.Fatalf("destination for %q = %q, want %q", category, got, want)
		}
	}
}

func TestDefaultConfigTemplateMatchesRepoExample(t *testing.T) {
	examplePath := filepath.Join("..", "..", "configs", "default.yaml")
	data, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatalf("reading %s: %v", examplePath, err)
	}

	if string(data) != DefaultConfigTemplate() {
		t.Fatalf("configs/default.yaml is out of sync with built-in default template")
	}
}
