package main

import (
	"fmt"
	"os"

	"github.com/theia-log/selene/cli"
)

func main() {
	selene := cli.NewCommandCLI("selene", "Client and agent for Theia log server.").
		AddCommand("watch", cli.WatcherCommand, "Watch for file changes. Runs selene in agent mode.").
		AddCommand("query", cli.QueryCommand, "Query the server for past and live events.").
		AddCommand("event", cli.EventCommand, "Generate event and publish to Theia server.")

	if err := selene.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "selene: %s\n", err.Error())
		os.Exit(1)
	}
}
