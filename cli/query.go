package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/theia-log/selene/comm"
)

func QueryCommand(args []string) error {
	queryFlags, flags := SetupQueryFlags()
	if err := flags.Parse(args); err != nil {
		return err
	}
	return RunQuery(queryFlags)
}

func RunQuery(flags *QueryFlags) error {
	serverURL, err := flags.GetServerURL()
	if err != nil {
		return err
	}
	client := comm.NewWebsocketClient(serverURL)
	filter, err := toQueryFilter(flags)
	if err != nil {
		return err
	}

	var resp chan *comm.EventResponse

	if flags.Live != nil && *flags.Live {
		resp, err = client.Receive(filter)
	} else {
		resp, err = client.Find(filter)
	}

	if err != nil {
		return err
	}

	done := make(chan bool)

	go func() {
		for {
			event, ok := <-resp
			if !ok {
				done <- true
				break
			}
			if event.Error != nil {
				fmt.Fprintln(os.Stderr, event.Error.Error())
				continue
			}
			// TODO: Format and print event
			ev := event.Event
			tags := ""
			if ev.Tags != nil {
				tags = strings.Join(ev.Tags, ",")
			}
			fmt.Printf("%s[%f](%s)%s: %s\n", ev.ID, ev.Timestamp, tags, ev.Source, ev.Content)
		}
	}()

	<-done
	return nil
}

func toQueryFilter(flags *QueryFlags) (*comm.EventFilter, error) {
	order := valueOrNil(flags.Order)
	if order != nil && *order != "asc" && *order != "desc" {
		return nil, fmt.Errorf("invalid value for order: %s", *order)
	}
	start := float64(0.0)
	if flags.Start != nil {
		start = *flags.Start
	}
	var end *float64
	if flags.End != nil && *flags.End != 0 {
		end = flags.End
	}
	filter := &comm.EventFilter{
		Content: valueOrNil(flags.Content),
		Order:   order,
		Tags:    flags.Tags,
		Start:   start,
		End:     end,
	}
	return filter, nil
}

func valueOrNil(str *string) *string {
	if str == nil {
		return nil
	}
	if *str == "" {
		return nil
	}
	return str
}
