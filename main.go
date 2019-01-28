package main

import (
	"flag"
	"fmt"

	"github.com/theia-log/selene/cli"
)

func main() {
	args := cli.SetupGlobalFlags()
	flag.Parse()
	fmt.Printf("args=%v\n", args)
}
