package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/theia-log/selene/model"

	"github.com/theia-log/selene/comm"
	"github.com/theia-log/selene/watcher"
)

// RunWatcher runs a watcher with the given watcher flags.
// It creates a WatcherDaemon and attaches the file, given as an input flag, as
// an EventSource.
// A connection to theia '/event' endpoint is created and the source events are
// pushed to the server.
func RunWatcher(args *WatcherFlags) error {
	daemon := watcher.NewWatchDaemon()

	if args.File == nil || *args.File == "" {
		return fmt.Errorf("no file to watch")
	}

	source, err := daemon.AddSource(*args.File, watcher.NewFileSource(*args.File))
	if err != nil {
		return err
	}

	serverURL, err := args.GetServerURL()
	if err != nil {
		return err
	}
	client := comm.NewWebsocketClient(serverURL)

	tags := []string{}
	if args.Tags != nil && len(args.Tags) > 0 {
		tags = args.Tags
	}

	source.OnSourceEvent(func(src string, diff []byte) {
		ev := &model.Event{
			ID:        uuid.Must(uuid.NewV4()).String(),
			Source:    src,
			Tags:      tags,
			Timestamp: float64(time.Now().UnixNano()) / float64(time.Millisecond),
			Content:   string(diff),
		}
		if err := client.Send(ev); err != nil {
			log.Println("Failed to send event ", err.Error())
		}
	})

	if err := daemon.Start(); err != nil {
		return err
	}

	done := make(chan bool)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP)
	go func() {
		for range c {
			if err := daemon.Stop(); err != nil {
				log.Println("Failed to stop watch daemon: ", err.Error())
			}
			done <- true
		}
	}()

	<-done

	return nil
}

// WatcherCommand implements the 'watch' subcommand.
// Takes a list of arguments to the watch subcommand, parses it and then calls
// RunWatcher with the parsed watcher flags.
func WatcherCommand(args []string) error {
	watcherFlags, flagSet := SetupWatcherFlags()
	err := flagSet.Parse(args)
	if err != nil {
		return err
	}
	return RunWatcher(watcherFlags)
}
