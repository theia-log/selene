package cli

import (
	"flag"
	"fmt"
	"strings"
)

type StringNVar []string

type GlobalFlags struct {
	ServerURL *string
	Verbose   *bool
	Host      *string
	Port      *int
}

type QueryFlags struct {
}

type WatcherFlags struct {
	*GlobalFlags
	File *string
	Tags StringNVar
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

func SetupWatcherFlags() (*WatcherFlags, *flag.FlagSet) {
	flags := flag.NewFlagSet("watch", flag.ExitOnError)
	watcherFlags := &WatcherFlags{
		GlobalFlags: SetupGlobalFlagsOn(flags),
		Tags:        StringNVar{},
	}
	watcherFlags.File = flags.String("f", "", "File to watch for changes")
	flags.Var(&watcherFlags.Tags, "t", "Tag the event.")
	return watcherFlags, flags
}

func SetupGlobalFlagsOn(fg *flag.FlagSet) *GlobalFlags {
	gf := &GlobalFlags{}

	gf.ServerURL = fg.String("server", "ws://localhost:6433", "Theia server URL.")
	gf.Verbose = fg.Bool("v", false, "Verbose output")
	gf.Host = fg.String("H", "", "Theia host")
	gf.Port = fg.Int("p", 0, "Theia port")

	return gf
}

func (gf *GlobalFlags) GetServerURL() (string, error) {
	if gf.ServerURL != nil {
		return *gf.ServerURL, nil
	}
	if gf.Host == nil {
		return "", fmt.Errorf("hostname missing")
	}

	if gf.Port == nil {
		return "", fmt.Errorf("port missing")
	}

	return fmt.Sprintf("ws://%s:%d", *gf.Host, *gf.Port), nil // FIXME: This assumes unsecure (ws) connection
}

func (n *StringNVar) String() string {
	if n == nil || *n == nil {
		return ""
	}
	return strings.Join(*n, ",")
}

func (n *StringNVar) Set(value string) error {
	*n = append(*n, value)
	return nil
}
