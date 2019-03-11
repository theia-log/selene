package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/theia-log/selene/comm"
)

func TestWatcherCommand(t *testing.T) {

	tempfile, err := ioutil.TempFile("", "watched")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tempfile.Name())

	fmt.Println(tempfile.Name())

	mock := comm.NewWebsocketMock().
		HandleReceivedMessage(func(data []byte) error {
			message := string(data)
			for _, line := range strings.Split(message, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "event:") {
					if len(line) == 6 {
						return fmt.Errorf("event preamble is malformed")
					}
				} else if strings.HasPrefix(line, "id:") {
					if len(line) == 3 {
						return fmt.Errorf("event ID is malformed")
					}
				} else if strings.HasPrefix(line, "timestamp:") {
					if len(line) == 10 {
						return fmt.Errorf("timestamp is malformed")
					}
				} else if strings.HasPrefix(line, "source:") {
					if line != fmt.Sprintf("source:%s", tempfile.Name()) {
						return fmt.Errorf("source value not set properly")
					}
				} else if strings.HasPrefix(line, "tags:") {
					if line != "tags:tag1,tag2" {
						return fmt.Errorf("tags value not set properly")
					}
				} else {
					if line != "test data" {
						return fmt.Errorf("event data is not correct - got: %s", line)
					}
				}
			}
			return nil
		})

	done := make(chan bool)

	go func() {
		mock.WaitRequestsToComplete(1)
		mock.Terminate()
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		done <- true
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		tempfile.Write([]byte("test data"))
		tempfile.Sync()
	}()

	if err := WatcherCommand([]string{
		"-server", mock.MockURL,
		"-f", tempfile.Name(),
		"-t", "tag1",
		"-t", "tag2",
	}); err != nil {
		t.Fatal(err)
	}

	<-done

	if mock.Errors != nil {
		for _, err := range mock.Errors {
			t.Log(err)
		}
		t.Fail()
	}
}
