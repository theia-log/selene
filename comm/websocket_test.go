package comm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
}

func (w *websocketMock) Expect(mesage string) *websocketMock {
	w.expect = []byte(mesage)
	return w
}

func (w *websocketMock) Respond(message string) *websocketMock {
	w.respond = append(w.respond, []byte(message))
	return w
}

func (w websocketMock) upgradedHandler(resp http.ResponseWriter, req *http.Request) {
	conn, err := w.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		http.Error(resp, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
	}
	mt, p, err := conn.ReadMessage()
	if err != nil {
		w.Errors = append(w.Errors, err)
		return
	}

	if w.expect != nil {
		if !w.expect.EqualsTo(p) {
			w.Errors = append(w.Errors, fmt.Errorf("expected '%s' but got '%s'", string(w.expect), string(p)))
			return
		}
	}

	if w.respond != nil {
		for _, msg := range w.respond {
			if err = conn.WriteMessage(mt, msg); err != nil {
				w.Errors = append(w.Errors, err)
			}
		}
	}

}

func NewWebsocketMock() *websocketMock {
	mock := &websocketMock{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(mock.upgradedHandler))
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"

	mock.MockURL = u.String()

	return mock
}

func TestWebsocketMock(t *testing.T) {
	mock := NewWebsocketMock()

	client := NewWebsocketClient(mock.MockURL)

	if err := client.Send(&model.Event{}); err != nil {
		t.Fatal(err)
	}

	if mock.Errors != nil {
		for _, err := range mock.Errors {
			t.Log(err)
		}
		t.Fail()
	}
}
