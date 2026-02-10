package cmd

import (
	"fmt"

	"github.com/geraldbahati/file-organizer/internal/daemon"
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
			fmt.Printf("  - %s\n", d)
		}
		fmt.Printf("Rules: %d categories\n", len(cfg.Rules))
		fmt.Printf("Debounce: %dms\n", cfg.DebounceMs)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
