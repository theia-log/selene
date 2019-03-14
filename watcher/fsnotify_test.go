package watcher

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestWatchForFileChanges(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	daemon := NewWatchDaemon()

	src, err := daemon.AddSource("test", NewFileSource(tmpFile.Name()))
	if err != nil {
		t.Fatal(err)
	}
	handled := false
	done := make(chan bool)
	src.OnSourceEvent(func(source string, diff []byte) {
		handled = true
		done <- true
		if source != tmpFile.Name() {
			t.Fatal("The source argument is passed incorrectly.")
		}
		if string(diff) != "test content" {
			t.Fatal("The content is passed incorrectly.")
		}

	})

	if err = daemon.Start(); err != nil {
		t.Fatal(err)
	}

	defer daemon.Stop()

	if _, err = tmpFile.Write([]byte("test content")); err != nil {
		t.Fatal(err)
	}

	if err = tmpFile.Sync(); err != nil {
		t.Fatal(err)
	}

	<-done
	if !handled {
		t.Fatal("Event not handled.")
	}
}
