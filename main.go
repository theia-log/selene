package main

import (
	"fmt"
	"os"

	"github.com/theia-log/selene/cli"
)

func main() {
	selene := cli.NewCommandCLI("selene", "Client and agent for Theia log server.").
		AddCommand("watch", cli.WatcherCommand, "Watch for file changes. Runs selene in agent mode.").
		AddCommand("query", cli.QueryCommand, "Query the server for past and live events.")

	if err := selene.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "selene: %s\n", err.Error())
		os.Exit(1)
	}
}
