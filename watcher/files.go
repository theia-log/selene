package watcher

type EventHandler func(source string, diff []byte)

type EventSource interface {
	OnSourceEvent(handler EventHandler)
	Trigger(diff []byte)
}

type WatchDaemon interface {
	Start() error
	Stop() error
	AddSource(source string, eventSource EventSource) (EventSource, error)
	RemoveSource(source string) error
}

type GenericEventSource struct {
	FilePath string
	handlers []EventHandler
}

func (g *GenericEventSource) OnSourceEvent(handler EventHandler) {
	g.handlers = append(g.handlers, handler)
}

func (g *GenericEventSource) Trigger(diff []byte) {
	for _, handler := range g.handlers {
		handler(g.FilePath, diff)
	}
}

func NewEventSource(filePath string) *GenericEventSource {
	return &GenericEventSource{
		FilePath: filePath,
		handlers: []EventHandler{},
	}
}
