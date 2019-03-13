package watcher

import "testing"

func TestAddEventHandler(t *testing.T) {
	eventSource := NewEventSource("test")

	eventSource.OnSourceEvent(func(source string, diff []byte) {})

	if eventSource.handlers == nil || len(eventSource.handlers) == 0 {
		t.Fatal("Expected the event handler to be registered.")
	}
}

func TestTriggerEventHandler(t *testing.T) {
	eventSource := NewEventSource("test")

	handlerCalled := false

	eventSource.OnSourceEvent(func(source string, diff []byte) {
		handlerCalled = true
		if source != "test" {
			t.Fatal("Invalid source passed to the handler function.")
		}
		if string(diff) != "trigger content" {
			t.Fatal("Invalid diff passed to the handler function.")
		}
	})

	eventSource.Trigger([]byte("trigger content"))

	if !handlerCalled {
		t.Fatal("Handler was not called.")
	}
}
