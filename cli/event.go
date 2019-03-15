package cli

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/theia-log/selene/comm"

	uuid "github.com/satori/go.uuid"
	"github.com/theia-log/selene/model"
)

func EventCommand(args []string) error {
	eventFlags, flagSet := SetupEventGeneratorFlags()
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	return RunEventGenerator(eventFlags)
}

func RunEventGenerator(flags *EventFlags) error {
	eventTemplate := &model.Event{
		ID:      asString(flags.ID),
		Source:  asString(flags.Source),
		Content: asString(flags.Content),
		Tags:    flags.Tags,
	}

	serverURL, err := flags.GetServerURL()
	if err != nil {
		return err
	}

	client := comm.NewWebsocketClient(serverURL)

	if flags.FromStdin != nil && (*flags.FromStdin) == true {
		return readFromStdinAndSend(eventTemplate, flags, client)
	}

	return sendOneAndExit(eventTemplate, flags, client)
}

func sendOneAndExit(template *model.Event, flags *EventFlags, client comm.Client) error {
	event := newFromTemplate(template)
	event.Content = ""
	if flags.Content != nil {
		event.Content = *flags.Content
	}
	return client.Send(event)
}

func readFromStdinAndSend(template *model.Event, flags *EventFlags, client comm.Client) error {
	reader := bufio.NewReader(os.Stdin)

	sep := "\n"
	if flags.Separator != nil && (*flags.Separator) != "" {
		sep = *flags.Separator
	}

	for {
		content, err := reader.ReadString(sep[0])
		eof := false
		if err != nil {
			if err.Error() == "EOF" {
				eof = true
			} else {
				return err
			}
		}

		if eof && content == "" {
			break
		}

		if flags.EofSeparator != nil && *flags.EofSeparator != "" {
			if idx := strings.Index(content, *flags.EofSeparator); idx >= 0 {
				content = content[0:idx]
				eof = true
			}
		}

		event := newFromTemplate(template)
		event.Content = content
		if err = client.Send(event); err != nil {
			return err
		}
		if eof {
			break
		}
	}

	return nil
}

func asString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func getTimestamp(timeVal *string) float64 {
	now := float64(time.Now().UnixNano()) / float64(time.Millisecond)
	if timeVal == nil {
		return now
	}
	timeStr := strings.TrimSpace(strings.ToLower(*timeVal))

	if timeStr == "now" {
		return now
	}

	return now
}

func newFromTemplate(template *model.Event) *model.Event {
	ev := &model.Event{
		ID:        template.ID,
		Source:    template.Source,
		Tags:      template.Tags,
		Timestamp: template.Timestamp,
	}

	if ev.ID == "" {
		ev.ID = uuid.Must(uuid.NewV4()).String()
	}

	if ev.Timestamp == 0.0 {
		ev.Timestamp = float64(time.Now().UnixNano()) / float64(time.Millisecond)
	}

	return ev
}
