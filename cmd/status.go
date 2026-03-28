package cmd

import (
	"fmt"

	"github.com/geraldbahati/file-organizer/internal/daemon"
	"github.com/geraldbahati/file-organizer/internal/fileutil"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the file organizer daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := daemon.ReadPID()
		if err != nil {
			return fmt.Errorf("reading PID: %w", err)
		}

		if pid > 0 && daemon.IsRunning(pid) {
			fmt.Printf("file-organizer is running (PID %d)\n", pid)
		} else {
			fmt.Println("file-organizer is not running")
		}

		fmt.Printf("\nConfig: %s\n", cfgPath)
		fmt.Printf("Watch dirs:\n")
		for _, d := range cfg.WatchDirs {
			expanded, err := fileutil.ExpandPath(d)
			if err != nil {
				return fmt.Errorf("expanding watch dir %q: %w", d, err)
			}
			fmt.Printf("  - %s\n", expanded)
		}
		fmt.Printf("Rules: %d categories\n", len(cfg.Rules))
		for _, rule := range cfg.Rules {
			expanded, err := fileutil.ExpandPath(rule.Destination)
			if err != nil {
				return fmt.Errorf("expanding destination for %q: %w", rule.Category, err)
			}
			fmt.Printf("  - %-20s %s\n", rule.Category, expanded)
		}
		fmt.Printf("Debounce: %dms\n", cfg.DebounceMs)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
