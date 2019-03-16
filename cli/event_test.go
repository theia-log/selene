package cli

import (
	"regexp"
	"testing"
)

func TestPatterns_rfc3339(t *testing.T) {
	pattern, err := regexp.Compile(rfc3339Pattern)
	if err != nil {
		panic(err)
	}

	if !pattern.Match([]byte("2019-10-11T12:13:14.900-0000")) {
		t.Fatal("Should match the time string.")
	}

	if !pattern.Match([]byte("2019-10-11T12:13:14.987+1234")) {
		t.Fatal("Should match the time string.")
	}

	if pattern.Match([]byte("2019-10-11")) {
		t.Fatal("Should not match the time string.")
	}
}

func TestPatterns_timestampPattern(t *testing.T) {
	pattern, err := regexp.Compile(timestampPattern)
	if err != nil {
		panic(err)
	}

	if !pattern.Match([]byte("123663713")) {
		t.Fatal("Expected to match the timestamp")
	}

	if !pattern.Match([]byte("123663713.909000")) {
		t.Fatal("Expected to match the timestamp as float")
	}

	if pattern.Match([]byte("123663713.1.111")) {
		t.Fatal("Expected NOT to match the invalid timestamp")
	}
}

func TestPatterns_manualTimePattern(t *testing.T) {
	pattern, err := regexp.Compile(manualTimePattern)
	if err != nil {
		panic(err)
	}

	if !pattern.Match([]byte("+3hrs")) {
		t.Fatal("Expected to match the manual time string.")
	}

	if !pattern.Match([]byte("-4000")) {
		t.Fatal("Expected to match the manual time string.")
	}

	if !pattern.Match([]byte("-23days")) {
		t.Fatal("Expected to match the manual time string.")
	}

	if pattern.Match([]byte("invalid")) {
		t.Fatal("Expected NOT to match the invalid manual time string.")
	}
}
