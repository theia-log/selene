package comm

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/theia-log/selene/model"

	"github.com/gorilla/websocket"
)

type mockTransport struct {
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("error:test")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
	}
	mt, p, err := conn.ReadMessage()
	if err != nil {
		log.Printf("cannot read message: %v", err)
		return
	}
	conn.WriteMessage(mt, []byte("hello "+string(p)))
}
func TestWebsocketMock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(websocketHandler))
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	http.DefaultTransport = &mockTransport{}

	client := NewWebsocketClient(u.String())

	if err := client.Send(&model.Event{}); err != nil {
		t.Fatal(err)
	}
}
