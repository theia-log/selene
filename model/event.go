package model

import (
	"fmt"
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
