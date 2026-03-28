package cmd

import (
	"log/slog"
	"os"

	"github.com/geraldbahati/file-organizer/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgPath string
	verbose bool
	cfg     config.Config
	logger  *slog.Logger
)

var rootCmd = &cobra.Command{
	Use:   "file-organizer",
	Short: "Automatically organize files on macOS",
	Long:  "A macOS background daemon that watches Downloads and Desktop, then moves files into predictable folders based on file type.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Set up logger.
		level := slog.LevelInfo
		if verbose {
			level = slog.LevelDebug
		}
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		}))

		// Skip config loading for version command.
		if cmd.Name() == "version" {
			return nil
		}

		// Determine config path.
		if cfgPath == "" {
			defaultPath, err := config.DefaultConfigPath()
			if err != nil {
				return err
			}
			cfgPath = defaultPath
		}

		var err error
		cfg, err = config.Load(cfgPath)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file path (default ~/.config/file-organizer/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
