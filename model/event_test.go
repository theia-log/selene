package model

import (
	"strings"
	"testing"
)

// TestLoad test loading of event from a string representation.
func TestLoad(t *testing.T) {
	ev := &Event{}
	err := ev.Load(strings.Join([]string{
		"event: 155 133 22",
		"id:331c531d-6eb4-4fb5-84d3-ea6937b01fdd",
		"timestamp: 1509989630.6749051",
		"source:/dev/sensors/door1-sensor",
		"tags:sensors,home,doors,door1",
		"Door has been unlocked",
	}, "\n"))
	if err != nil {
		t.Fatal(err)
	}

	if ev.ID != "331c531d-6eb4-4fb5-84d3-ea6937b01fdd" {
		t.Fatal("ID was not parsed correctly.")
	}
	if !(ev.Timestamp < 1509989630.67491 && ev.Timestamp > 1509989630.67490) {
		t.Fatal("Timestamp not parsed correctly. Got: ", ev.Timestamp)
	}
	if ev.Source != "/dev/sensors/door1-sensor" {
		t.Fatal("Source was not parsed correctly.")
	}
	if ev.Content != "Door has been unlocked" {
		t.Fatal("Content was not parsed correctly.")
	}

	if ev.Tags == nil {
		t.Fatal("Tags were not parsed at all.")
	}

	if len(ev.Tags) != 4 {
		t.Fatal("Tags were not parsed correctly.")
	}

	expected := []string{"sensors", "home", "doors", "door1"}
	for i := 0; i < 4; i++ {
		if ev.Tags[i] != expected[i] {
			t.Fatalf("Tags parsed incorrectly. Tags[%d] is %s but expected to be %s.", i, ev.Tags[i], expected[i])
		}
	}
}

// TestLoadBytes tests loading of event data from serialized event as bytes.
func TestLoadBytes(t *testing.T) {
	eventBytes := []byte(strings.Join([]string{
		"event: 155 133 22",
		"id:331c531d-6eb4-4fb5-84d3-ea6937b01fdd",
		"timestamp: 1509989630.6749051",
		"source:/dev/sensors/door1-sensor",
		"tags:sensors,home,doors,door1",
		"Door has been unlocked",
	}, "\n"))
	ev := &Event{}
	err := ev.LoadBytes(eventBytes)
	if err != nil {
		t.Fatal(err)
	}

	if ev.ID != "331c531d-6eb4-4fb5-84d3-ea6937b01fdd" {
		t.Fatal("ID was not parsed correctly.")
	}
	if !(ev.Timestamp < 1509989630.67491 && ev.Timestamp > 1509989630.67490) {
		t.Fatal("Timestamp not parsed correctly. Got: ", ev.Timestamp)
	}
	if ev.Source != "/dev/sensors/door1-sensor" {
		t.Fatal("Source was not parsed correctly.")
	}
	if ev.Content != "Door has been unlocked" {
		t.Fatal("Content was not parsed correctly.")
	}

	if ev.Tags == nil {
		t.Fatal("Tags were not parsed at all.")
	}

	if len(ev.Tags) != 4 {
		t.Fatal("Tags were not parsed correctly.")
	}

	expected := []string{"sensors", "home", "doors", "door1"}
	for i := 0; i < 4; i++ {
		if ev.Tags[i] != expected[i] {
			t.Fatalf("Tags parsed incorrectly. Tags[%d] is %s but expected to be %s.", i, ev.Tags[i], expected[i])
		}
	}
}

// TestDump tests dumping of an event as string.
func TestDump(t *testing.T) {
	expected := strings.Join([]string{
		"event:153 131 22",
		"id:331c531d-6eb4-4fb5-84d3-ea6937b01fdd",
		"timestamp:1509989630.674905",
		"source:/dev/sensors/door1-sensor",
		"tags:sensors,home,doors,door1",
		"Door has been unlocked",
	}, "\n")
	ev := &Event{
		ID:        "331c531d-6eb4-4fb5-84d3-ea6937b01fdd",
		Source:    "/dev/sensors/door1-sensor",
		Timestamp: 1509989630.6749051,
		Tags:      []string{"sensors", "home", "doors", "door1"},
		Content:   "Door has been unlocked",
	}
	data, err := ev.Dump()
	if err != nil {
		t.Fatal(err)
	}
	if data != expected {
		t.Fatalf("Event not dumped properly.\nExpected:\n>\n%s\n>\nBut instead got:\n<\n%s\n<\n", expected, data)
	}
}

// TestDumpBytes tests dumping of event as a slice of bytes.
func TestDumpBytes(t *testing.T) {
	expected := strings.Join([]string{
		"event:153 131 22",
		"id:331c531d-6eb4-4fb5-84d3-ea6937b01fdd",
		"timestamp:1509989630.674905",
		"source:/dev/sensors/door1-sensor",
		"tags:sensors,home,doors,door1",
		"Door has been unlocked",
	}, "\n")
	ev := &Event{
		ID:        "331c531d-6eb4-4fb5-84d3-ea6937b01fdd",
		Source:    "/dev/sensors/door1-sensor",
		Timestamp: 1509989630.6749051,
		Tags:      []string{"sensors", "home", "doors", "door1"},
		Content:   "Door has been unlocked",
	}
	data, err := ev.DumpBytes()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != expected {
		t.Fatalf("Event not dumped properly.\nExpected:\n>\n%s\n>\nBut instead got:\n<\n%s\n<\n", expected, string(data))
	}
}
