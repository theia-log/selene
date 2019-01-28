package cli

import "flag"

type GlobalFlags struct {
	ServerURL *string
	Verbose   *bool
	Host      *string
	Port      *int
}

func SetupGlobalFlags() *GlobalFlags {
	gf := &GlobalFlags{}

	gf.ServerURL = flag.String("server", "ws://localhost:6433", "Theia server URL.")
	gf.Verbose = flag.Bool("v", false, "Verbose output")

	return gf
}
