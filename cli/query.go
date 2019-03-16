package cli

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/theia-log/selene/model"

	"github.com/theia-log/selene/comm"
)

// QueryCommand implements the 'query' subcommand.
// Takes a list of arguments to the query subcommand, parses it and then calls
// RunQuery with the parsed query flags.
func QueryCommand(args []string) error {
	queryFlags, flags := SetupQueryFlags()
	if err := flags.Parse(args); err != nil {
		return err
	}
	return RunQuery(queryFlags)
}

// RunQuery runs a query against the server with the given query flags.
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

	colors := NewAuroraColors()
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
			ev := event.Event
			PrintEvent(ev, DefaultEventFormat, colors)
		}
	}()

	<-done
	return nil
}

// toQueryFilter transforms the QueryFlags to an EventFilter ready to be passed
// down to the theia client.
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

	var filterOrder *comm.EventOrder
	if order != nil {
		ord := comm.EventOrder(*order)
		filterOrder = &ord
	}

	filter := &comm.EventFilter{
		Content: valueOrNil(flags.Content),
		Order:   filterOrder,
		Tags:    flags.Tags,
		Start:   start,
		End:     end,
	}
	return filter, nil
}

// valueOrNil returns nil if the passed pointer is nil or points to an empty
// string (""). Otherwise returns the original string.
func valueOrNil(str *string) *string {
	if str == nil {
		return nil
	}
	if *str == "" {
		return nil
	}
	return str
}

type templateEvent struct {
	ID        string
	IDShort   string
	Timestamp string
	Tags      string
	Source    string
	Content   string
}

// FullEventFormat format template for printing an Event - full data.
var FullEventFormat = "{{ .ID }}:[{{ .Timestamp }}]({{ .Source }}) {{ .Tags }} - {{ .Content }}"

// ShortEventFormat format template to print event content, source and tags - short format.
var ShortEventFormat = "[{{ .Source }}]{{ .Tags }} - {{ .Content }}"

// DefaultEventFormat default format template for printing an event.
var DefaultEventFormat = "{{ .IDShort }}:[{{ .Timestamp }}]({{ .Source }}) {{ .Tags }} - {{ .Content }}"

// PrintEvent prints the event using the provided format template to STDOUT.
func PrintEvent(event *model.Event, format string, colors Colors) {
	content := event.Content
	if !strings.HasSuffix(content, "\n") {
		content = content + "\n"
	}
	te := &templateEvent{
		ID:        colors.ColoredText(Context{"color": "secondary"}, event.ID),
		IDShort:   colors.ColoredText(Context{"color": "secondary"}, fmt.Sprintf("%s", event.ID[0:7])),
		Content:   colors.ColoredContent(Context{}, content),
		Source:    colors.ColoredText(Context{"color": "secondary"}, event.Source),
		Timestamp: colors.ColoredText(Context{"color": "info"}, fmt.Sprintf("%f", event.Timestamp)),
	}

	tags := []string{}
	if event.Tags != nil {
		for _, tag := range event.Tags {
			tags = append(tags, colors.ColoredTag(Context{}, tag))
		}
	}
	te.Tags = strings.Join(tags, " ")

	tpl, err := template.New("event").Parse(format)
	if err != nil {
		panic(err)
	}

	if err = tpl.Execute(os.Stdout, te); err != nil {
		panic(err)
	}
}
