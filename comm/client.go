package comm

import "github.com/theia-log/selene/model"

type EventFilter struct {
}

func (f *EventFilter) DumpBytes() ([]byte, error) {
	return nil, nil
}

type EventResponse struct {
	Event *model.Event
	Error error
}

type Client interface {
	Send(event model.Event) error
	Receive(filter *EventFilter) (chan *EventResponse, error)
	Find(filter *EventFilter) (chan *EventResponse, error)
}
