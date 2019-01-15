package watcher

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FSNotifyEventSource struct {
	*GenericEventSource
}

func (s *FSNotifyEventSource) AttachToWatcher(watcher *fsnotify.Watcher) error {
	return watcher.Add(s.FilePath)
}

func (s *FSNotifyEventSource) DetachFromWatcher(watcher *fsnotify.Watcher) error {
	return watcher.Remove(s.FilePath)
}

type FSNotifyWatcher struct {
	watcher *fsnotify.Watcher
	sources map[string]EventSource
	started bool
	mux     sync.Mutex
}

func (f *FSNotifyWatcher) Start() error {
	f.mux.Lock()
	defer f.mux.Unlock()
	if f.started {
		return fmt.Errorf("already started")
	}
	f.started = true
	return nil
}

func (f *FSNotifyWatcher) Stop() error {
	f.mux.Lock()
	defer f.mux.Unlock()
	if !f.started {
		return fmt.Errorf("stopped")
	}
	f.watcher.Close()
	return nil
}

func (f *FSNotifyWatcher) AddSource(source string, eventSource EventSource) (EventSource, error) {
	f.mux.Lock()
	src, ok := f.sources[source]
	if ok {
		f.mux.Unlock()
		return src, nil
	}
	if fsnSource, ok := eventSource.(*FSNotifyEventSource); ok {
		if err := fsnSource.AttachToWatcher(f.watcher); err != nil {
			f.mux.Unlock()
			return nil, err
		}
	}
	f.sources[source] = eventSource
	f.mux.Unlock()

	return eventSource, nil
}

func (f *FSNotifyWatcher) RemoveSource(source string) error {
	f.mux.Lock()
	src, ok := f.sources[source]
	if !ok {
		f.mux.Unlock()
		return fmt.Errorf("source not managed")
	}
	delete(f.sources, source)
	f.mux.Unlock()

	if fsnSource, ok := src.(*FSNotifyEventSource); ok {
		// try to remove from watcher
		if err := fsnSource.DetachFromWatcher(f.watcher); err != nil {
			return err
		}
	}

	// TODO: generic dispose call

	return nil
}
