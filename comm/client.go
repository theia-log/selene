package comm

import (
	"encoding/json"

	"github.com/theia-log/selene/model"
)

type EventFilter struct {
	Start   float64  `json:"start,omitempty"`
	End     *float64 `json:"end,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Content *string  `json:"content,omitempty"`
	Order   *string  `json:"order,omitempty"`
}

func (f *EventFilter) DumpBytes() ([]byte, error) {
	return json.Marshal(f)
}

type EventResponse struct {
	Event *model.Event
	Error error
}

type Client interface {
	Send(event *model.Event) error
	Receive(filter *EventFilter) (chan *EventResponse, error)
	Find(filter *EventFilter) (chan *EventResponse, error)
}
