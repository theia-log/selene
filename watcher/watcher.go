package watcher

// EventHandler defines a type for handling a change in a particular source.
// The function takes two parameters: a source, the event source name, and
// the difference, array of bytes of the actual change.
// This is designed to handle changes that primarily append to the event source,
// so the default interpretation is that the diff is the value that have been
// appended to the previous value of the source:
//	 State(event_now) - State(event_before)
type EventHandler func(source string, diff []byte)

// EventSource defines the general interface for an event source.
type EventSource interface {
	// OnSourceEvent adds an EventHandler for the generated events by this
	// source.
	OnSourceEvent(handler EventHandler)

	// Trigger triggers an event on this source, passing the diff value.
	// Note that this function, although exposed, it is mainly used internally
	// or by EventSource managers that can keep track of the underlying sources
	// like watching for file changes or reacting to appends in the system log.
	Trigger(diff []byte)
}

// WatchDaemon is a general interface for manager of event source.
type WatchDaemon interface {
	// Start the daemon. Events from the event sources shall be handled after
	// this point.
	Start() error

	// Stop the daemon. No events will be handled after this call and the daemon
	// will try to close all managed EventSources.
	Stop() error

	// AddSource adds an EventSource to be managed by this WatchDaemon.
	// Each source is identified by its name. This name can later be used to
	// remove the source event from this daemon.
	AddSource(source string, eventSource EventSource) (EventSource, error)

	// Removes the source event from this daemon. The underlying source for the
	// source event shall not be managed after this point.
	RemoveSource(source string) error
}

// GenericEventSource implements the abstract and common functionalities of an
// EventSource.
type GenericEventSource struct {
	// FilePath is the source name identified by a file path.
	FilePath string

	// List of event handlers to be triggered on source event.
	handlers []EventHandler
}

// OnSourceEvent registers an EventHandler to be triggered when this source is
// changed (produces an event).
func (g *GenericEventSource) OnSourceEvent(handler EventHandler) {
	g.handlers = append(g.handlers, handler)
}

// Trigger an event over this source. Used mostly internally.
func (g *GenericEventSource) Trigger(diff []byte) {
	for _, handler := range g.handlers {
		handler(g.FilePath, diff)
	}
}

// NewEventSource builds new GenericEventSource for the given file path.
func NewEventSource(filePath string) *GenericEventSource {
	return &GenericEventSource{
		FilePath: filePath,
		handlers: []EventHandler{},
	}
}
