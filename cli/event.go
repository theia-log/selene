package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/theia-log/selene/comm"

	uuid "github.com/satori/go.uuid"
	"github.com/theia-log/selene/model"
)

// EventCommand implements the 'event' subcommand.
// Takes a list of arguments to the event subcommand, parses it and then calls
// RunEventGenerator with the parsed query flags.
func EventCommand(args []string) error {
	eventFlags, flagSet := SetupEventGeneratorFlags()
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	return RunEventGenerator(eventFlags)
}

// RunEventGenerator generates new event with the given properties and sends it
// to the Theia server.
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

var rfc3339Pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}((-|\\+|Z)\\d{2}:\\d{2}){0,1}$"
var timestampPattern = "^\\d+(\\.\\d+){0,1}$"
var manualTimePattern = "^(\\+|-)\\d+(\\w{1,7}){0,1}$"

type timeStringParser func(timeStr string) (float64, error)

func rfc3339Parser(timeStr string) (float64, error) {
	tm, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0.0, err
	}
	return float64(tm.UnixNano()) / float64(time.Millisecond), nil
}

func timestampParser(timeStr string) (float64, error) {
	return strconv.ParseFloat(timeStr, 64)
}

var units = map[string]float64{
	"ms,millisecond,milliseconds": 1.0,
	"s,second,seconds":            1000.0,
	"m,min,minute,minutes":        60 * 1000.0,
	"h,hr,hrs,hour,hours":         60 * 60 * 1000.0,
	"d,day,days":                  24 * 60 * 60 * 1000.0,
	"w,week,weeks":                7 * 24 * 60 * 60 * 1000.0,
	"mn,mon,month,months":         30 * 24 * 60 * 60 * 1000.0,
	"y,yr,year,years":             365 * 24 * 60 * 60 * 1000.0,
}

func manualStringParser(timeStr string) (float64, error) {
	timeStr = strings.ToLower(timeStr)
	now := float64(time.Now().UnixNano()) / (float64(time.Millisecond))
	sign := timeStr[0]
	mul := 1.0

	if sign == '-' {
		mul = -1.0
	}

	unit, value := splitAlphaNum(timeStr[1:])

	timeVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, err
	}

	for unitsList, multiplier := range units {
		match := false
		aliases := strings.Split(unitsList, ",")
		for _, u := range aliases {
			if u == unit {
				mul *= multiplier
				match = true
				break
			}
		}
		if match {
			break
		}
	}

	timeVal *= mul

	return now + timeVal, nil
}

func splitAlphaNum(str string) (string, string) {
	alpha := ""
	num := ""
	for i := 0; i < len(str); i++ {
		if unicode.IsLetter(rune(str[i])) {
			alpha = str[i:len(str)]
			break
		} else {
			num = num + string(str[i])
		}
	}
	return alpha, num
}

var parsers = map[string]timeStringParser{
	rfc3339Pattern:    rfc3339Parser,
	timestampPattern:  timestampParser,
	manualTimePattern: manualStringParser,
}

func parseTime(timeStr string) (float64, error) {
	for pattern, parser := range parsers {
		match, err := regexp.MatchString(pattern, timeStr)
		if err != nil {
			return 0.0, err
		}
		if match {
			return parser(timeStr)
		}
	}
	return 0.0, fmt.Errorf("invalid time string")
}
