package model

import (
	"bufio"
	"bytes"
	"fmt"
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
	return
}

func (ev *Event) LoadBytes(eventData []byte) (err error) {
	reader := bytes.NewReader(eventData)
	_, headerSize, contentSize, err := parsePreamble(bufio.NewReader(reader))
	if err != nil {
		return err
	}

	buff := make([]byte, headerSize)
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
		// line := scanner.Text()

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

func parsePreamble(event *bufio.Reader) (total, header, content int64, err error) {
	lnbytes, err := event.ReadBytes('\n')
	if err != nil {
		return 0, 0, 0, err
	}
	line := string(lnbytes)
	if !strings.HasPrefix(line, "event:") {
		return 0, 0, 0, fmt.Errorf("invalid preamble")
	}
	parts := strings.Split(strings.TrimSpace(line)[6:], " ")
	total, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		header, err = strconv.ParseInt(parts[1], 10, 64)
	}
	if err != nil {
		content, err = strconv.ParseInt(parts[2], 10, 64)
	}
	return
}
