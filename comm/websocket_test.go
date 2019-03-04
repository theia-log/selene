package comm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/theia-log/selene/model"

	"github.com/gorilla/websocket"
)

type Message []byte

func (m Message) EqualsTo(message []byte) bool {
	return string(m) == string(message)
}

type websocketMock struct {
	expect   Message
	respond  []Message
	MockURL  string
	upgrader websocket.Upgrader
	Errors   []error
	done     chan bool
}

func (w *websocketMock) Expect(mesage string) *websocketMock {
	w.expect = []byte(mesage)
	return w
}

func (w *websocketMock) Respond(message string) *websocketMock {
	w.respond = append(w.respond, []byte(message))
	return w
}

func (w *websocketMock) AddError(err error) *websocketMock {
	if w.Errors == nil {
		w.Errors = []error{}
	}
	w.Errors = append(w.Errors, err)
	return w
}

func (w *websocketMock) markRequestCompleted() {
	w.done <- true
}

func (w *websocketMock) WaitRequestsToComplete(n int) {
	for ; n > 0; n-- {
		<-w.done
	}
}

func (w *websocketMock) upgradedHandler(resp http.ResponseWriter, req *http.Request) {
	conn, err := w.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		http.Error(resp, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
		w.markRequestCompleted()
		return
	}
	mt, p, err := conn.ReadMessage()
	if err != nil {
		w.AddError(err)
		w.markRequestCompleted()
		return
	}
	if w.expect != nil {
		if !w.expect.EqualsTo(p) {
			w.AddError(fmt.Errorf("expected '%s' but got '%s'", string(w.expect), string(p)))
			w.markRequestCompleted()
			return
		}
	}

	if w.respond != nil {
		for _, msg := range w.respond {
			if err = conn.WriteMessage(mt, msg); err != nil {
				w.AddError(err)
				w.markRequestCompleted()
				return
			}
		}
	}
	w.markRequestCompleted()
}

func NewWebsocketMock() *websocketMock {
	mock := &websocketMock{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		done: make(chan bool, 5),
	}

	srv := httptest.NewServer(http.HandlerFunc(mock.upgradedHandler))
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"

	mock.MockURL = u.String()

	return mock
}

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
