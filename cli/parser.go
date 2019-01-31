package cli

import "flag"

type GlobalFlags struct {
	ServerURL *string
	Verbose   *bool
	Host      *string
	Port      *int
}

type QueryFlags struct {
}

type WatcherFlags struct {
}

func SetupGlobalFlags() *GlobalFlags {
	gf := &GlobalFlags{}

	gf.ServerURL = flag.String("server", "ws://localhost:6433", "Theia server URL.")
	gf.Verbose = flag.Bool("v", false, "Verbose output")
	gf.Host = flag.String("H", "", "Theia host")
	gf.Port = flag.Int("p", 0, "Theia port")

	return gf
}

func SetupQueryFlags() (*QueryFlags, *flag.FlagSet) {
	return nil, nil
}

func SetupWatcherFlags() (*QueryFlags, *flag.FlagSet) {
	return nil, nil
}

func SetupGlobalFlagsOn(fg *flag.FlagSet) *GlobalFlags {
	gf := &GlobalFlags{}

	gf.ServerURL = fg.String("server", "ws://localhost:6433", "Theia server URL.")
	gf.Verbose = fg.Bool("v", false, "Verbose output")
	gf.Host = fg.String("H", "", "Theia host")
	gf.Port = fg.Int("p", 0, "Theia port")

	return gf
}
