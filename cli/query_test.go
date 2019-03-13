package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/theia-log/selene/model"

	"github.com/theia-log/selene/comm"
)

func TestQueryCommand(t *testing.T) {
	mock := comm.NewWebsocketMock().
		Expect(strings.Join([]string{
			"{",
			"\"start\":100.1,",
			"\"end\":200.2,",
			"\"tags\":[\"tag1\",\"tag2\"],",
			"\"content\":\"content-match\"",
			"}",
		}, "")).
		Respond("ok")
	done := make(chan bool)
	go func() {
		mock.WaitRequestsToComplete(1)
		mock.Terminate()

		if mock.Errors != nil {
			for _, err := range mock.Errors {
				t.Log(err.Error())
			}
			t.FailNow()
		}
		done <- true
	}()
	err := QueryCommand([]string{"-server", mock.MockURL,
		"-s", "100.1",
		"-e", "200.2",
		"-c", "content-match",
		"-t", "tag1", "-t", "tag2",
	})
	if err != nil {
		t.Fatal(err)
	}
	<-done
}

func TestPrintEvent(t *testing.T) {
	colors := NewAuroraColors()
	ev := &model.Event{
		ID:        "684b84e9-14aa-4267-828a-55ea371c0508",
		Timestamp: float64(time.Now().UnixNano()) / float64(time.Millisecond),
		Source:    "/tmp/source",
		Tags:      []string{"tag1", "tag2", "error", "info", "add"},
		Content:   "This is the content.",
	}

	PrintEvent(ev, DefaultEventFormat, colors)

}
