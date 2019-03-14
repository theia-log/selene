package cli

import (
	"flag"
	"fmt"
	"strings"
)

// StringNVar implements the flag.Value interface for flags that can hold
// multiple values.
type StringNVar []string

// GlobalFlags holds the parsed values for the global flags. These flags are
// shared (common) between all subcommands like watch, query and event.
type GlobalFlags struct {
	// ServerURL is the full URL to the theia server. It is a websocket URL.
	ServerURL *string

	// Verbose is a flag for verbose output. If set, selene should provide
	// more verbose output, printing more messages to stdout.
	Verbose *bool

	// Host is the theia host. This is the domain part of the server URL.
	// Either host and port must be given, or a full ServerURL.
	Host *string

	// Port is the theia port. This is the port part of the server URL.
	// Either host and port must be given, or a full ServerURL.
	Port *int
}

// QueryFlags holds the parsed values for subcommand 'query'.
// This struct holds the global flags, and specific flags for the query command.
// Most of the values are filter values for matching events on the server.
type QueryFlags struct {
	*GlobalFlags

	// Start timestamp. Match events that occurred after or at this time.
	Start *float64

	// End timestamp. Match events that occurred before or at this time.
	End *float64

	// Tags list of tags to match. The values may be a regular expression.
	Tags StringNVar

	// Content is a regular expression to match the events content.
	Content *string

	// Order is the order in which the matched events should be returned. It
	// can be 'asc' - ascending, or 'desc' - descending order by the event
	// timestamp.
	Order *string

	// Live is a flag to indicate whether to query live (real-time) events.
	Live *bool
}

type EventFlags struct {
	ID           *string
	Source       *string
	Time         *string
	Tags         StringNVar
	Content      *string
	Separator    *string
	EofSeparator *string
	FromStdin    *bool
}

// WatcherFlags holds the parsed values for the subcommand 'watch'.
// This struct holds the global flags, and additional flags for the watch
// command.
type WatcherFlags struct {
	*GlobalFlags

	// File is the path of the file to be watched for changes.
	File *string

	// Tags is a list of tags to be attached to the generated events.
	Tags StringNVar
}

// SetupQueryFlags creates a FlagSet for parsing the 'query' subcommand and
// creates a wrapper QueryFlags to hold the parsed values from the command line.
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

// SetupWatcherFlags creates a FlagSet for parsing the 'watcher' subcommand and
// creates a wrapper WatcherFlags to hold the parsed values from the command
// line.
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

// SetupEventGeneratorFlags creates a FlagSet for parsing the 'event'
// subcommand.
func SetupEventGeneratorFlags() (*EventFlags, *flag.FlagSet) {
	return nil, nil
}

// SetupGlobalFlagsOn adds the global flags to an existing FlagSet and returns
// the wrapper struct that will hold the parsed values for the flags.
func SetupGlobalFlagsOn(fg *flag.FlagSet) *GlobalFlags {
	gf := &GlobalFlags{}

	gf.ServerURL = fg.String("server", "ws://localhost:6433", "Theia server URL.")
	gf.Verbose = fg.Bool("v", false, "Verbose output")
	gf.Host = fg.String("H", "", "Theia host")
	gf.Port = fg.Int("p", 0, "Theia port")

	return gf
}

// GetServerURL returns a valid server URL based on the set up flags.
// If ServerURL is set, then that value is preferred and returned.
// If not, the URL is generated from the Host and Port values. If any of them
// are not set, then an error is returned.
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

// String returns a string representation of the multiple values flag value.
func (n *StringNVar) String() string {
	if n == nil || *n == nil {
		return ""
	}
	return strings.Join(*n, ",")
}

// Set adds an additional value to the list of multiple flag values.
func (n *StringNVar) Set(value string) error {
	*n = append(*n, value)
	return nil
}
