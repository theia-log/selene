package comm

import "github.com/theia-log/selene/model"

type Client interface {
	Connect() error
	Send(event model.Event) error
	Receive() (chan model.Event, error)
}
