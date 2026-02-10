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

// Watcher wraps fsnotify with per-file debounce timers and a single worker
// goroutine that processes files via a buffered channel.
type Watcher struct {
	fsw      *fsnotify.Watcher
	handler  FileHandler
	debounce time.Duration
	logger   *slog.Logger
	mu       sync.Mutex
	timers   map[string]*time.Timer
	workCh   chan string
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
		workCh:   make(chan string, 64),
	}, nil
}

// Add adds a directory to watch.
func (w *Watcher) Add(dir string) error {
	return w.fsw.Add(dir)
}

// Run processes events until ctx is cancelled. Blocks.
func (w *Watcher) Run(ctx context.Context) error {
	defer w.fsw.Close()

	// Single worker goroutine drains the work channel.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for path := range w.workCh {
			w.handler(path)
		}
	}()

	var err error
	for {
		select {
		case <-ctx.Done():
			w.cancelAll()
			close(w.workCh)
			wg.Wait()
			err = ctx.Err()
			return err
		case event, ok := <-w.fsw.Events:
			if !ok {
				close(w.workCh)
				wg.Wait()
				return nil
			}
			if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Rename) == 0 {
				continue
			}
			w.debounceFile(event.Name)
		case fsErr, ok := <-w.fsw.Errors:
			if !ok {
				close(w.workCh)
				wg.Wait()
				return nil
			}
			w.logger.Error("watcher error", "error", fsErr)
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
		w.workCh <- path
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
