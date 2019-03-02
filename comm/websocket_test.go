package comm

import (
	"fmt"
	"net/http"
	"testing"
)

type mockTransport struct {
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("error:test")
}

func TestWebsocketMock(t *testing.T) {
	http.DefaultTransport = &mockTransport{}

	client := NewWebsocketClient("ws://localhost:11211")

	if err := client.Send(nil); err != nil {
		t.Fatal(err)
	}
}
