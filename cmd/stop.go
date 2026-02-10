package cmd

import (
	"fmt"

	"github.com/geraldbahati/file-organizer/internal/daemon"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running file organizer daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.StopDaemon(); err != nil {
			return err
		}
		fmt.Println("file-organizer stopped")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
