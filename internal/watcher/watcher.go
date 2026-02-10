package watcher

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileHandler is called when a file event fires after debounce.
type FileHandler func(path string)

// Watcher wraps fsnotify with per-file debounce timers.
type Watcher struct {
	fsw        *fsnotify.Watcher
	handler    FileHandler
	debounce   time.Duration
	logger     *slog.Logger
	mu         sync.Mutex
	timers     map[string]*time.Timer
}

// New creates a Watcher. debounce is how long to wait after the last event
// for a given file before calling handler.
func New(handler FileHandler, debounce time.Duration, logger *slog.Logger) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		fsw:      fsw,
		handler:  handler,
		debounce: debounce,
		logger:   logger,
		timers:   make(map[string]*time.Timer),
	}, nil
}

// Add adds a directory to watch.
func (w *Watcher) Add(dir string) error {
	return w.fsw.Add(dir)
}

// Run processes events until ctx is cancelled. Blocks.
func (w *Watcher) Run(ctx context.Context) error {
	defer w.fsw.Close()

	for {
		select {
		case <-ctx.Done():
			w.cancelAll()
			return ctx.Err()
		case event, ok := <-w.fsw.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Rename) == 0 {
				continue
			}
			w.debounceFile(event.Name)
		case err, ok := <-w.fsw.Errors:
			if !ok {
				return nil
			}
			w.logger.Error("watcher error", "error", err)
		}
	}
}

func (w *Watcher) debounceFile(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if t, exists := w.timers[path]; exists {
		t.Stop()
	}

	w.timers[path] = time.AfterFunc(w.debounce, func() {
		w.mu.Lock()
		delete(w.timers, path)
		w.mu.Unlock()
		w.handler(path)
	})
}

func (w *Watcher) cancelAll() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for path, t := range w.timers {
		t.Stop()
		delete(w.timers, path)
	}
}
