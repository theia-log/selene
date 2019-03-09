package cli

import (
	"flag"
	"testing"
)

func TestGetServerURL_GlobalFlags(t *testing.T) {
	host := "server"
	port := 9876
	serverURL := "ws://custom:9876"
	gf := &GlobalFlags{
		Host: &host,
		Port: &port,
	}

	url, err := gf.GetServerURL()
	if err != nil {
		t.Fatal(err)
	}

	if url != "ws://server:9876" {
		t.Fatalf("Expected 'ws://server:9876' but got '%s'\n", url)
	}

	gf = &GlobalFlags{
		ServerURL: &serverURL,
	}

	url, err = gf.GetServerURL()
	if err != nil {
		t.Fatal(err)
	}

	if url != "ws://custom:9876" {
		t.Fatalf("Expected 'ws://custom:9876' but got '%s'\n", url)
	}

	gf = &GlobalFlags{}

	if _, err = gf.GetServerURL(); err == nil {
		t.Fatal("Expected to fail when getting ServerURL from empty struct.")
	}

}

func TestSetupQueryFlags(t *testing.T) {
	qf, fs := SetupQueryFlags()

	if qf == nil {
		t.Fatal("Expected QueryFlags structure.")
	}

	if fs == nil {
		t.Fatal("Expected FlagsSet to be set.")
	}

	if err := fs.Parse([]string{"-s", "11.1", "-e", "12.2", "-c", "content",
		"-sort", "desc", "-t", "tag1", "-t", "tag2", "-live"}); err != nil {
		t.Fatal(err)
	}

	if qf.Start == nil || *qf.Start != 11.1 {
		t.Fatal("Start timestamp not parsed properly")
	}

	if qf.End == nil || *qf.End != 12.2 {
		t.Fatal("End timestamp not parsed properly")
	}

	if qf.Content == nil || *qf.Content != "content" {
		t.Fatal("Content flag not parsed properly")
	}

	if qf.Order == nil || *qf.Order != "desc" {
		t.Fatal("Order flag is not parsed properly")
	}

	if qf.Tags == nil || len(qf.Tags) != 2 {
		t.Fatal("Tags list not parsed properly")
	}

	if qf.Tags[0] != "tag1" || qf.Tags[1] != "tag2" {
		t.Fatal("Tags values not parsed properly")
	}

	if qf.Live == nil || *qf.Live != true {
		t.Fatal("Live flag not parsed")
	}
}

func TestGlobalFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ExitOnError)

	gf := SetupGlobalFlagsOn(fs)

	if err := fs.Parse([]string{"-H", "test-host", "-p", "9876", "-v",
		"-server", "ws://server:9090"}); err != nil {
		t.Fatal(err)
	}

	if gf.Host == nil || *gf.Host != "test-host" {
		t.Fatal("Host property not parsed properly")
	}

	if gf.Port == nil || *gf.Port != 9876 {
		t.Fatal("Port flag not parsed properly")
	}

	if gf.ServerURL == nil || *gf.ServerURL != "ws://server:9090" {
		t.Fatal("Server flag not parsed properly")
	}

	if gf.Verbose == nil || *gf.Verbose != true {
		t.Fatal("Verbose flag not parsed")
	}
}

func TestSetupWatcherFlags(t *testing.T) {
	wf, fs := SetupWatcherFlags()

	if wf == nil {
		t.Fatal("Expected to get WatcherFlags struct")
	}

	if fs == nil {
		t.Fatal("Expected to ge the flag set")
	}

	if err := fs.Parse([]string{"-f", "file1",
		"-t", "tag1", "-t", "tag2"}); err != nil {
		t.Fatal(err)
	}

	if wf.File == nil || *wf.File != "file1" {
		t.Fatal("File flag not parsed properly")
	}

	if wf.Tags == nil || len(wf.Tags) != 2 {
		t.Fatal("Tags flags not parsed properly")
	}

	if wf.Tags[0] != "tag1" || wf.Tags[1] != "tag2" {
		t.Fatal("Tags values not parsed properly")
	}
}
