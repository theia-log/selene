package cli

import (
	"github.com/theia-log/selene/comm"
	"github.com/theia-log/selene/model"
)

func RunCli(args []string) error {

	client := comm.NewWebsocketClient("ws://localhost:6433")

	if err := client.Send(&model.Event{
		ID: "test",
	}); err != nil {
		panic(err)
	}

	return nil
}
