// Package watcher defines a watcher daemon that handles events from diffrent
// event sources.
// This package defines the interfaces for the WatcherDaemon and EventSource.
// Implementations for file-based event sources and watcher based on inotify
// is provided as well.
//
// An example on watching a file for changes:
//	import (
//		"github.com/theia-log/selene/watcher"
//	)
//
//	// Tail-like watcher for a log file
//	func main() {
//		daemon := watcher.NewWatchDaemon()
//		src := watcher.NewFileSource("/var/log/nginx/error.log")
//
//		// Add handler to print the diff.
//		src.OnSourceEvent(func(file string, diff []byte) {
//			fmt.Printf("%s: %s\n", file, string(diff))
//		});
//
//		// Attach the event source to the daemon.
//		if _, err := daemon.AddSource(src); err != nil {
//			panic(err)
//		}
//	}
//
package watcher
