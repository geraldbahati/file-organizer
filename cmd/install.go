package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/geraldbahati/file-organizer/internal/daemon"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install macOS launchd agent (start on login)",
	RunE: func(cmd *cobra.Command, args []string) error {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("finding executable: %w", err)
		}
		exe, err = filepath.Abs(exe)
		if err != nil {
			return fmt.Errorf("resolving executable path: %w", err)
		}

		plistPath, err := daemon.LaunchdPlistPath()
		if err != nil {
			return err
		}

		plist := daemon.GeneratePlist(exe)

		if err := os.MkdirAll(filepath.Dir(plistPath), 0755); err != nil {
			return fmt.Errorf("creating LaunchAgents dir: %w", err)
		}

		if err := os.WriteFile(plistPath, []byte(plist), 0644); err != nil {
			return fmt.Errorf("writing plist: %w", err)
		}

		// Load the agent.
		if out, err := exec.Command("launchctl", "load", plistPath).CombinedOutput(); err != nil {
			return fmt.Errorf("launchctl load: %s: %w", string(out), err)
		}

		fmt.Printf("Installed launchd agent: %s\n", plistPath)
		fmt.Println("file-organizer will start automatically on login.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
