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
	*GlobalFlags
	Start   *float64
	End     *float64
	Tags    StringNVar
	Content *string
	Order   *string
	Live    *bool
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
	flags := flag.NewFlagSet("query", flag.ExitOnError)
	queryFlags := &QueryFlags{
		GlobalFlags: SetupGlobalFlagsOn(flags),
		Tags:        StringNVar{},
	}

	queryFlags.Start = flags.Float64("s", 0.0, "Start timestamp")
	queryFlags.End = flags.Float64("e", 0.0, "End timestamp")
	queryFlags.Content = flags.String("c", "", "Match event content (regular expression)")
	queryFlags.Order = flags.String("sort", "", "Sort order (only if not live). Possible values are asc or desc.")
	queryFlags.Live = flags.Bool("live", false, "Whether to query for live events in real time.")

	flags.Var(&queryFlags.Tags, "t", "Match if any tag with this value (regular expression).")

	return queryFlags, flags
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
