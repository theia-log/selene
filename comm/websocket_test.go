package comm

import (
	"strings"
	"testing"

	"github.com/theia-log/selene/model"
)

// TestWebsocketClientSend tests vanilla case of sending an Event via the
// WebsocketClient API.
func TestWebsocketClientSend(t *testing.T) {
	mock := NewWebsocketMock().Expect(strings.Join([]string{
		"event:71 65 6",
		"id:id-001",
		"timestamp:1551733035.230000",
		"source:/src",
		"tags:tag1,tag2",
		"event1",
	}, "\n"))

	client := NewWebsocketClient(mock.MockURL)

	if err := client.Send(&model.Event{
		ID:        "id-001",
		Source:    "/src",
		Timestamp: 1551733035.23,
		Tags:      []string{"tag1", "tag2"},
		Content:   "event1",
	}); err != nil {
		t.Fatal(err)
	}

	mock.WaitRequestsToComplete(1)
	if mock.Errors != nil {
		for _, err := range mock.Errors {
			t.Log(err)
		}
		t.Fail()
	}
}

func TestWebsocketClientReceive(t *testing.T) {
	mock := NewWebsocketMock().
		Expect("{\"start\":10}").
		Respond("ok").
		Respond(strings.Join([]string{
			"event:71 65 6",
			"id:id-001",
			"timestamp:1551733035.230000",
			"source:/src",
			"tags:tag1,tag2",
			"event1",
		}, "\n"))

	client := NewWebsocketClient(mock.MockURL)

	resp, err := client.Receive(Filter(10.0))

	if err != nil {
		t.Fatal(err)
	}

	mock.WaitRequestsToComplete(1)

	if mock.Errors != nil {
		for _, err := range mock.Errors {
			t.Log(err)
		}
		t.Fail()
		return
	}

	event := <-resp
	if event.Error != nil {
		t.Fatal(event.Error)
	}

	if event.Event == nil {
		t.Fatal("Expected to get a parsed event.")
	}

	if event.Event.ID != "id-001" {
		t.Fatal("Event not parsed properly")
	}
}

func TestWebsocketClientFind(t *testing.T) {
	mock := NewWebsocketMock().
		Expect("{\"start\":10,\"content\":\"event1\"}").
		Respond("ok").
		Respond(strings.Join([]string{
			"event:71 65 6",
			"id:id-001",
			"timestamp:1551733035.230000",
			"source:/src",
			"tags:tag1,tag2",
			"event1",
		}, "\n"))

	client := NewWebsocketClient(mock.MockURL)

	resp, err := client.Find(Filter(10.0).MatchContent("event1"))

	if err != nil {
		t.Fatal(err)
	}

	mock.WaitRequestsToComplete(1)

	if mock.Errors != nil {
		for _, err := range mock.Errors {
			t.Log(err)
		}
		t.Fail()
		return
	}

	event := <-resp
	if event.Error != nil {
		t.Fatal(event.Error)
	}

	if event.Event == nil {
		t.Fatal("Expected to get a parsed event.")
	}

	if event.Event.ID != "id-001" {
		t.Fatal("Event not parsed properly")
	}
}
