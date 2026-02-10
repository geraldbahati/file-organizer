package config

// Rule maps a category to its destination and file extensions.
type Rule struct {
	Category    string   `yaml:"category"`
	Extensions  []string `yaml:"extensions"`
	Destination string   `yaml:"destination"`
}

// DefaultRules returns the built-in file categorization rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Category:    "Images",
			Extensions:  []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico", ".tiff", ".heic"},
			Destination: "~/Pictures",
		},
		{
			Category:    "Videos",
			Extensions:  []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm"},
			Destination: "~/Movies",
		},
		{
			Category:    "Music",
			Extensions:  []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma", ".m4a"},
			Destination: "~/Music",
		},
		{
			Category:    "Documents",
			Extensions:  []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".rtf", ".odt", ".csv"},
			Destination: "~/Documents",
		},
		{
			Category:    "Archives",
			Extensions:  []string{".zip", ".tar", ".gz", ".rar", ".7z", ".bz2", ".xz"},
			Destination: "~/Documents/Archives",
		},
		{
			Category:    "Code",
			Extensions:  []string{".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".rs", ".rb", ".html", ".css", ".json", ".yaml", ".yml", ".xml", ".sh", ".sql"},
			Destination: "~/Documents/Code",
		},
		{
			Category:    "Apps",
			Extensions:  []string{".dmg", ".pkg"},
			Destination: "~/Downloads/Apps",
		},
		{
			Category:    "Fonts",
			Extensions:  []string{".ttf", ".otf", ".woff", ".woff2"},
			Destination: "~/Documents/Fonts",
		},
		{
			Category:    "eBooks",
			Extensions:  []string{".epub", ".mobi"},
			Destination: "~/Documents/Books",
		},
		{
			Category:    "Torrents",
			Extensions:  []string{".torrent"},
			Destination: "~/Documents/Torrents",
		},
		{
			Category:    "Design",
			Extensions:  []string{".psd", ".ai", ".sketch", ".fig", ".xd"},
			Destination: "~/Documents/Design",
		},
	}
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		WatchDirs:  []string{"~/Downloads", "~/Desktop"},
		Rules:      DefaultRules(),
		DebounceMs: 500,
	}
}
