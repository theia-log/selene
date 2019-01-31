package main

import (
	"flag"

	"github.com/theia-log/selene/cli"
)

func main() {
	cli.SetupGlobalFlags()
	flag.Parse()
}
