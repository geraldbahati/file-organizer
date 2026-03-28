package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Rule maps a category to its destination and file extensions.
type Rule struct {
	Category    string   `yaml:"category"`
	Extensions  []string `yaml:"extensions"`
	Destination string   `yaml:"destination"`
}

const defaultConfigTemplate = `# file-organizer configuration
# macOS defaults:
# - Downloads is treated as an inbox, not a permanent dumping ground.
# - Archives and installers stay inside Downloads so they do not feel "lost."
# - Media and working files move into the standard macOS home folders.
# Example: a new .zip downloaded to ~/Downloads will be moved to ~/Downloads/Archives.

watch_dirs:
  # Primary inbox for files you fetch from browsers, Slack, AirDrop, and mail.
  - ~/Downloads
  # Keeps the Desktop clear without changing the default macOS save locations.
  - ~/Desktop

debounce_ms: 500

rules:
  - category: Images
    extensions: [.jpg, .jpeg, .png, .gif, .bmp, .svg, .webp, .ico, .tiff, .heic, .heif, .icns]
    destination: ~/Pictures/Imports

  - category: Videos
    extensions: [.mp4, .mkv, .avi, .mov, .wmv, .flv, .webm, .m4v]
    destination: ~/Movies/Imports

  - category: Music
    extensions: [.mp3, .wav, .flac, .aac, .ogg, .wma, .m4a, .aif, .aiff]
    destination: ~/Music/Imports

  - category: PDFs
    extensions: [.pdf]
    destination: ~/Documents/PDFs

  - category: Writing
    extensions: [.doc, .docx, .txt, .rtf, .odt, .pages, .md]
    destination: ~/Documents/Writing

  - category: Spreadsheets
    extensions: [.xls, .xlsx, .csv, .numbers]
    destination: ~/Documents/Spreadsheets

  - category: Presentations
    extensions: [.ppt, .pptx, .key]
    destination: ~/Documents/Presentations

  - category: Archives
    extensions: [.zip, .tar, .gz, .rar, .7z, .bz2, .xz]
    destination: ~/Downloads/Archives

  - category: Disk Images
    extensions: [.dmg, .iso, .img]
    destination: ~/Downloads/Disk Images

  - category: Installer Packages
    extensions: [.pkg, .xip]
    destination: ~/Downloads/Installers

  - category: Code
    extensions: [.go, .py, .js, .ts, .java, .c, .cpp, .rs, .rb, .html, .css, .json, .yaml, .yml, .xml, .sh, .sql]
    destination: ~/Documents/Code

  - category: Fonts
    extensions: [.ttf, .otf, .woff, .woff2]
    destination: ~/Downloads/Fonts

  - category: eBooks
    extensions: [.epub, .mobi]
    destination: ~/Documents/Books

  - category: Torrents
    extensions: [.torrent]
    destination: ~/Downloads/Torrents

  - category: Design
    extensions: [.psd, .ai, .sketch, .fig, .xd]
    destination: ~/Documents/Design
`

// DefaultConfigTemplate returns the built-in config template written for new installs.
func DefaultConfigTemplate() string {
	return defaultConfigTemplate
}

// DefaultRules returns the built-in file categorization rules.
func DefaultRules() []Rule {
	return DefaultConfig().Rules
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	var cfg Config
	if err := yaml.Unmarshal([]byte(defaultConfigTemplate), &cfg); err != nil {
		panic(fmt.Sprintf("invalid built-in default config: %v", err))
	}
	return cfg
}
