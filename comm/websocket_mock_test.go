package comm

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestBuildWebsocketMock(t *testing.T) {
	mock := NewWebsocketMock().Expect("test").Respond("response")

	if mock.MockURL == "" {
		t.Fatal("Expected to generate mock url for the test websocket server.")
	}
}

func TestWebsocketMockCall(t *testing.T) {
	mock := NewWebsocketMock().Expect("request").Respond("response")

	c, _, err := websocket.DefaultDialer.Dial(mock.MockURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.WriteMessage(websocket.TextMessage, []byte("request")); err != nil {
		t.Fatal(err)
	}

	_, msg, err := c.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if msg == nil {
		t.Fatal("Expected to get a message")
	}

	if string(msg) != "response" {
		t.Fatal("Got unexpected response from mock server")
	}
}
