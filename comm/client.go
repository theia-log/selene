package comm

import (
	"encoding/json"

	"github.com/theia-log/selene/model"
)

// EventOrder is the order in which the events should be returned. Can be either
// 'asc' - ascending, or 'desc' - descending.
type EventOrder string

// OrderAsc sort in ascending order.
const OrderAsc EventOrder = "asc"

// OrderDesc sort in descending order.
const OrderDesc EventOrder = "desc"

// EventFilter holds values for filtering events.
// This structure is used when filtering past events and filering real-time
// events as well.
type EventFilter struct {
	// Match events that happened after this time. This is required for
	// filtering both past and real-time events.
	Start float64 `json:"start,omitempty"`

	// Match events that happened before this timestamp. Optional.
	End *float64 `json:"end,omitempty"`

	// Tags is a list of possible values to match for tags. Each value may be a
	// regular expression. Matches the event only if all patterns are found in
	// the event tag list.
	Tags []string `json:"tags,omitempty"`

	// Match the content of the event. This value is evaluated as regular
	// expression.
	Content *string `json:"content,omitempty"`

	// Order in which to return the events. Makes sense only for past events.
	Order *EventOrder `json:"order,omitempty"`
}

// DumpBytes serializes the event filter values as bytes.
// Theia expects the filter in JSON format, so this function serializes the
// filter data to JSON, then encodes in UTF-8.
func (f *EventFilter) DumpBytes() ([]byte, error) {
	return json.Marshal(f)
}

// MatchContent sets the matcher for the content of this EventFilter.
// Returns pointer to this EventFilter.
func (f *EventFilter) MatchContent(content string) *EventFilter {
	f.Content = &content
	return f
}

// MatchTag adds a matcher to the list of tags matchers of this EventFilter.
// Returns pointer to this EventFilter.
func (f *EventFilter) MatchTag(tag ...string) *EventFilter {
	f.Tags = append(f.Tags, tag...)
	return f
}

// MatchEnd sets the end timestamp for this EventFilter. Match events that
// happened before this time.
// Returns pointer to this EventFilter.
func (f *EventFilter) MatchEnd(end float64) *EventFilter {
	f.End = &end
	return f
}

// OrderAsc set the filter order to ascending.
// Returns pointer to this EventFilter.
func (f *EventFilter) OrderAsc() *EventFilter {
	order := OrderAsc
	f.Order = &order
	return f
}

// OrderDesc sets the filter order to descending.
// Returns pointer to this EventFilter.
func (f *EventFilter) OrderDesc() *EventFilter {
	order := OrderDesc
	f.Order = &order
	return f
}

// Filter creates new EventFilter with start timestamp.
func Filter(start float64) *EventFilter {
	return &EventFilter{
		Start: start,
	}
}

// EventResponse holds an Event or an Error.
// Used to pass data over a channel from a query operation.
type EventResponse struct {
	// The event received from the server. In case of error it may be nil.
	Event *model.Event

	// The error that occurred.
	// If the read was successful, this will be set to nil.
	Error error
}

// Client interface describes the client API for a theia server.
// Defines methods to send events and to query both past and real-time events.
// The querying operations (Receive, Find) are both streaming and asynchronous,
// meaning that the number of events to be returned is not known and the server
// returns (streams) the events as the arrive.
// These functions return a chan to listen on, and events are decoded and
// published on the channel as they arrive.
type Client interface {
	// Send publishes an event to the remote server.
	// Returns an error if the client fails to send the event.
	Send(event *model.Event) error

	// Receive opens a channel for real-time events to the server.
	// The events that match the EventFilter are returned of the EventResponse
	// channel.
	// If the client fails to open a real-time event channel to the server, an
	// error is returned.
	// It should be noted that the server will never close this type of channel.
	// The responsibility for closing the connection is on the client side.
	Receive(filter *EventFilter) (chan *EventResponse, error)

	// Find performs a lookup for past events on the server.
	// The server will return all the events that match the EventFilter.
	// The events are returned as they are found and are published on the
	// EventResponse channel.
	// If the client fails to open the channel or other error occurs during
	// establishing the connection or while setting the filter, an error will
	// be returned.
	// The server will automatically close the connection once all of the
	// matching events have been returned to the client.
	Find(filter *EventFilter) (chan *EventResponse, error)
}
