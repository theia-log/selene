package comm

import "github.com/theia-log/selene/model"

type EventFilter struct {
}

type Client interface {
	Send(event model.Event) error
	Receive(filter EventFilter) (chan model.Event, error)
	Find(filter EventFilter) (chan model.Event, error)
}
