package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/theia-log/selene/model"

	"github.com/theia-log/selene/comm"
	"github.com/theia-log/selene/watcher"
)

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

	source.OnSourceEvent(func(src string, diff []byte) {
		ev := &model.Event{
			ID:        uuid.Must(uuid.NewV4()).String(),
			Source:    src,
			Tags:      []string{},
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
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			// TODO: Stop client here
			done <- true
		}
	}()

	<-done

	return nil
}

func WatcherCommand(args []string) error {
	watcherFlags, flagSet := SetupWatcherFlags()
	err := flagSet.Parse(args)
	if err != nil {
		return err
	}
	return RunWatcher(watcherFlags)
}
