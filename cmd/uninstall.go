package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/geraldbahati/file-organizer/internal/daemon"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove macOS launchd agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		plistPath, err := daemon.LaunchdPlistPath()
		if err != nil {
			return err
		}

		if _, err := os.Stat(plistPath); os.IsNotExist(err) {
			fmt.Println("launchd agent is not installed")
			return nil
		}

		// Unload the agent (ignore errors if not loaded).
		exec.Command("launchctl", "unload", plistPath).Run()

		if err := os.Remove(plistPath); err != nil {
			return fmt.Errorf("removing plist: %w", err)
		}

		fmt.Printf("Removed launchd agent: %s\n", plistPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
