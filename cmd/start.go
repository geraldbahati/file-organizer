package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/geraldbahati/file-organizer/internal/daemon"
	"github.com/geraldbahati/file-organizer/internal/organizer"
	"github.com/geraldbahati/file-organizer/internal/watcher"
	"github.com/spf13/cobra"
)

var foreground bool

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the file organizer daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		if foreground {
			return runForeground()
		}
		return runBackground()
	},
}

func init() {
	startCmd.Flags().BoolVar(&foreground, "foreground", false, "run in foreground (don't daemonize)")
	rootCmd.AddCommand(startCmd)
}

func runForeground() error {
	// Check if already running.
	pid, _ := daemon.ReadPID()
	if pid > 0 && daemon.IsRunning(pid) && pid != os.Getpid() {
		return fmt.Errorf("daemon already running (PID %d)", pid)
	}

	if err := daemon.WritePID(); err != nil {
		return fmt.Errorf("writing PID file: %w", err)
	}
	defer daemon.RemovePID()

	index, err := organizer.NewRuleIndex(cfg.Rules)
	if err != nil {
		return fmt.Errorf("building rule index: %w", err)
	}
	org := organizer.New(index, logger)

	debounce := time.Duration(cfg.DebounceMs) * time.Millisecond
	w, err := watcher.New(func(path string) {
		if _, err := org.ProcessFile(path); err != nil {
			logger.Error("failed to process file", "path", path, "error", err)
		}
	}, debounce, logger)
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}

	dirs, err := cfg.ExpandedWatchDirs()
	if err != nil {
		return fmt.Errorf("expanding watch dirs: %w", err)
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			logger.Warn("watch directory does not exist, skipping", "dir", dir)
			continue
		}
		if err := w.Add(dir); err != nil {
			return fmt.Errorf("watching %q: %w", dir, err)
		}
		logger.Info("watching directory", "dir", dir)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		logger.Info("received signal, shutting down", "signal", sig)
		cancel()
	}()

	logger.Info("file-organizer started", "pid", os.Getpid())
	return w.Run(ctx)
}

func runBackground() error {
	pid, _ := daemon.ReadPID()
	if pid > 0 && daemon.IsRunning(pid) {
		return fmt.Errorf("daemon already running (PID %d)", pid)
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding executable: %w", err)
	}

	cmdArgs := []string{"start", "--foreground"}
	if cfgPath != "" {
		cmdArgs = append(cmdArgs, "--config", cfgPath)
	}
	if verbose {
		cmdArgs = append(cmdArgs, "--verbose")
	}

	proc := exec.Command(exe, cmdArgs...)
	proc.Stdout = nil
	proc.Stderr = nil
	proc.Stdin = nil

	if err := proc.Start(); err != nil {
		return fmt.Errorf("starting background process: %w", err)
	}

	fmt.Printf("file-organizer started in background (PID %d)\n", proc.Process.Pid)
	return nil
}
