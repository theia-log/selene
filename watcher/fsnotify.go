package watcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FSNotifyEventSource struct {
	*GenericEventSource
	AbsFilePath string
	ParentDir   string
	currentPos  int64
}

func (s *FSNotifyEventSource) AttachToWatcher(watcher *fsnotify.Watcher) error {
	return watcher.Add(s.ParentDir)
}

func (s *FSNotifyEventSource) DetachFromWatcher(watcher *fsnotify.Watcher) error {
	return watcher.Remove(s.ParentDir)
}

func (s *FSNotifyEventSource) fileAvailable() error {
	f, err := os.Open(s.AbsFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	endPos, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		return err
	}
	s.currentPos = endPos
	return nil
}

func (s *FSNotifyEventSource) calculateDiff() ([]byte, error) {
	f, err := os.Open(s.AbsFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	curr := s.currentPos

	_, err = f.Seek(curr, os.SEEK_SET)
	if err != nil {
		return nil, err
	}

	diff, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	curr, err = f.Seek(0, os.SEEK_CUR)
	if err != nil {
		return nil, err
	}

	s.currentPos = curr

	return diff, err
}

func (s *FSNotifyEventSource) handleFSNotifyEvent(ev fsnotify.Event) error {
	if ev.Op&fsnotify.Write == fsnotify.Write {
		diff, err := s.calculateDiff()
		if err != nil {
			return err
		}
		s.Trigger(diff)
	} else if ev.Op&fsnotify.Create == fsnotify.Create {
		if err := s.fileAvailable(); err != nil {
			return err
		}
		diff, err := s.calculateDiff()
		if err != nil {
			return err
		}
		s.Trigger(diff)
	}
	return nil
}

func NewFileSource(filePath string) EventSource {
	parentDir, absPath, err := GetFileParts(filePath)
	if err != nil {
		panic(err) // FIXME
	}

	fileSource := &FSNotifyEventSource{
		GenericEventSource: NewEventSource(filePath),
		AbsFilePath:        absPath,
		ParentDir:          parentDir,
	}

	fileSource.fileAvailable()

	return fileSource
}

type FSNotifyWatcher struct {
	watcher     *fsnotify.Watcher
	watchedDirs map[string][]EventSource
	sources     map[string]EventSource
	started     bool
	mux         sync.Mutex
	done        chan bool
}

func (f *FSNotifyWatcher) Start() error {
	f.mux.Lock()
	defer f.mux.Unlock()
	if f.started {
		return fmt.Errorf("already started")
	}
	f.started = true
	f.listenForChanges()
	return nil
}

func (f *FSNotifyWatcher) Stop() error {
	f.mux.Lock()
	defer f.mux.Unlock()
	if !f.started {
		return fmt.Errorf("stopped")
	}
	f.watcher.Close()
	<-f.done
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
		sourcesInDir, ok := f.watchedDirs[fsnSource.ParentDir]
		if !ok {
			sourcesInDir = []EventSource{}
			if err := fsnSource.AttachToWatcher(f.watcher); err != nil {
				f.mux.Unlock()
				return nil, err
			}
		}
		sourcesInDir = append(sourcesInDir, fsnSource)
		f.watchedDirs[fsnSource.ParentDir] = sourcesInDir
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

func (f *FSNotifyWatcher) listenForChanges() {
	go func() {
		defer func() {
			f.done <- true
		}()
		for {
			select {
			case event, ok := <-f.watcher.Events:
				if !ok {
					return
				}
				f.handleWatcherEvent(event)
			case err, ok := <-f.watcher.Errors:
				if !ok {
					return
				}
				log.Println("[ERR]: Watcher event error: ", err.Error())
			}
		}
	}()
}

func (f *FSNotifyWatcher) handleWatcherEvent(ev fsnotify.Event) {
	parentDir, absPath, err := GetFileParts(ev.Name)
	if err != nil {
		log.Println("error: cannot determine abs path for event: ", ev.Name)
		return
	}
	sources, ok := f.watchedDirs[parentDir]
	if !ok {
		return
	}

	for _, source := range sources {
		fsource := source.(*FSNotifyEventSource)
		if fsource.AbsFilePath == absPath {
			go func() {
				if err = fsource.handleFSNotifyEvent(ev); err != nil {
					log.Println("Error in handling event: ", err.Error())
				}
			}()
		}
	}
}

func NewWatchDaemon() WatchDaemon {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err) // FIME:
	}
	return &FSNotifyWatcher{
		mux:         sync.Mutex{},
		sources:     map[string]EventSource{},
		watchedDirs: map[string][]EventSource{},
		started:     false,
		watcher:     fsWatcher,
		done:        make(chan bool),
	}
}

func GetFileParts(filePath string) (parentDir string, absPath string, err error) {
	absPath, err = filepath.Abs(filePath)
	if err != nil {
		return
	}
	parentDir = filepath.Dir(absPath)
	return
}
