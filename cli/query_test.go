package cli

import (
	"strings"
	"testing"

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
