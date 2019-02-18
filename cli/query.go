package cli

import (
	"fmt"

	"github.com/theia-log/selene/comm"
)

func QueryCommand(args []string) error {
	return nil
}

func RunQuery(flags *QueryFlags) error {
	serverURL, err := flags.GetServerURL()
	if err != nil {
		return err
	}
	client := comm.NewWebsocketClient(serverURL)

	resp, err := client.Find(&comm.EventFilter{})
	if err != nil {
		return err
	}

	go func() {
		for {
			event, ok := <-resp
			if !ok {
				break
			}
			fmt.Println(event.Event)
		}
	}()
	return nil
}
