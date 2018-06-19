package watcher

type EventHandler func(source string, diff []byte)

type EventSource interface {
	OnSourceEvent(handler EventHandler)
}

type WatchDaemon interface {
	Start() error
	Stop() error
	AddSource(source string, eventSource EventSource)
	RemoveSource(source string)
}
