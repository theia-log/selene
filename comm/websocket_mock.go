package comm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gorilla/websocket"
)

// Message is a websocket message represented as an array of bytes.
type Message []byte

// OnMessageHandler gets called when a message is received.
type OnMessageHandler func([]byte) error

// EqualsTo check if this message is equal to another message data.
func (m Message) EqualsTo(message []byte) bool {
	return string(m) == string(message)
}

// WebsocketMock implements a mock specification for websocket server.
type WebsocketMock struct {
	expect         Message
	requestHandler OnMessageHandler
	respond        []Message
	MockURL        string
	upgrader       websocket.Upgrader
	Errors         []error
	done           chan bool
	conn           *websocket.Conn
}

// Expect expect to receive a message with the given value.
func (w *WebsocketMock) Expect(mesage string) *WebsocketMock {
	w.expect = []byte(mesage)
	return w
}

// Respond responds to the client websocket with the given message after the
// first message has been received.
func (w *WebsocketMock) Respond(message string) *WebsocketMock {
	w.respond = append(w.respond, []byte(message))
	return w
}

// AddError adds an error to the mock object. The errors are kept sequentially
// as they are added.
func (w *WebsocketMock) AddError(err error) *WebsocketMock {
	if w.Errors == nil {
		w.Errors = []error{}
	}
	w.Errors = append(w.Errors, err)
	return w
}

func (w *WebsocketMock) markRequestCompleted() {
	w.done <- true
}

// WaitRequestsToComplete can be called to wait until the whole handling of
// incoming and outgoing messages has been completed. You must specify the
// number of messages to be handled before the execution can continue and this
// method returns control.
func (w *WebsocketMock) WaitRequestsToComplete(n int) {
	if w.Errors != nil || len(w.Errors) > 0 {
		return
	}
	for ; n > 0; n-- {
		<-w.done
	}
}

// upgradedHandler handles the receiving and responding to websocket messages
// using httptest package structures and Gorilla/websocket upgrader for the
// standard HTTP test server.
func (w *WebsocketMock) upgradedHandler(resp http.ResponseWriter, req *http.Request) {
	conn, err := w.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		http.Error(resp, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
		w.markRequestCompleted()
		return
	}
	w.conn = conn
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
	if w.requestHandler != nil {
		if err = w.requestHandler(p); err != nil {
			w.AddError(err)
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

// HandleReceivedMessage add handler for received messages.
func (w *WebsocketMock) HandleReceivedMessage(handler OnMessageHandler) *WebsocketMock {
	w.requestHandler = handler
	return w
}

// Terminate terminates and closes the server connection.
func (w *WebsocketMock) Terminate() error {
	if w.conn != nil {
		return w.conn.Close()
	}
	return fmt.Errorf("no connection")
}

// NewWebsocketMock constructs a new websocket mock to be used when testing.
func NewWebsocketMock() *WebsocketMock {
	mock := &WebsocketMock{
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
