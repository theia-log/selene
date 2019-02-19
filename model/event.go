package model

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Event struct {
	ID        string
	Timestamp float64
	Source    string
	Tags      []string
	Content   string
}

func (ev *Event) Load(eventData string) (err error) {
	return ev.LoadBytes([]byte(eventData))
}

func (ev *Event) LoadBytes(eventData []byte) (err error) {
	reader := bytes.NewReader(eventData)
	preamble, _, headerSize, contentSize, err := parsePreamble(bufio.NewReader(reader))
	if err != nil {
		return err
	}

	buff := make([]byte, headerSize)
	_, err = reader.Seek(preamble, io.SeekStart)
	if err != nil {
		return err
	}
	if read, err := reader.Read(buff); err != nil || int64(read) != headerSize {
		if err != nil {
			return err
		}
		return fmt.Errorf("corrupted header")
	}
	scanner := bufio.NewScanner(bytes.NewReader(buff))

	for {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line != "" {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			switch key {
			case "id":
				ev.ID = value
			case "source":
				ev.Source = value
			case "timestamp":
				if ev.Timestamp, err = strconv.ParseFloat(value, 64); err != nil {
					return err
				}
			case "tags":
				ev.Tags = strings.Split(value, ",")
			default:
				// ignore unknown key
			}
		}

	}
	buff = make([]byte, contentSize)
	if read, err := reader.Read(buff); err != nil || int64(read) != contentSize {
		if err != nil {
			return err
		}
		return fmt.Errorf("corrupted content")
	}

	ev.Content = string(buff)

	return
}

func (ev *Event) Dump() (eventData string, err error) {
	event := ev.dump()
	contentSize := len([]byte(ev.Content))
	totalSize := len([]byte(event))
	headerSize := totalSize - contentSize
	preamble := fmt.Sprintf("event:%d %d %d", totalSize, headerSize, contentSize)
	return fmt.Sprintf("%s\n%s", preamble, event), nil
}

func (ev *Event) dump() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("id:%s\n", ev.ID))
	builder.WriteString(fmt.Sprintf("timestamp:%f\n", ev.Timestamp))
	builder.WriteString(fmt.Sprintf("source:%s\n", ev.Source))
	if ev.Tags != nil {
		builder.WriteString(fmt.Sprintf("tags:%s\n", strings.Join(ev.Tags, ",")))
	}

	builder.WriteString(ev.Content)

	return builder.String()
}

func (ev *Event) DumpBytes() (eventData []byte, err error) {
	evString, err := ev.Dump()
	eventData = []byte(evString)
	return
}

func parsePreamble(event *bufio.Reader) (preamble, total, header, content int64, err error) {
	lnbytes, err := event.ReadBytes('\n')
	if err != nil {
		return 0, 0, 0, 0, err
	}
	preamble = int64(len(lnbytes))
	line := string(lnbytes)
	if !strings.HasPrefix(line, "event:") {
		return 0, 0, 0, 0, fmt.Errorf("invalid preamble")
	}
	parts := strings.Split(strings.TrimSpace(strings.TrimSpace(line)[6:]), " ")
	total, err = strconv.ParseInt(parts[0], 10, 64)
	if err == nil {
		header, err = strconv.ParseInt(parts[1], 10, 64)
	}
	if err == nil {
		content, err = strconv.ParseInt(parts[2], 10, 64)
	}
	return
}
