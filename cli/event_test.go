package cli

import (
	"math"
	"regexp"
	"testing"
	"time"
)

func currentTimeMillis() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Millisecond)
}

func assertEqual(expected, actual, tolerance float64, t *testing.T) {
	if math.Abs(expected-actual) > tolerance {
		t.Fatalf("Expected %f but got %f (tol %f)", expected, actual, tolerance)
	}
}

func TestPatterns_rfc3339(t *testing.T) {
	pattern, err := regexp.Compile(rfc3339Pattern)
	if err != nil {
		panic(err)
	}

	if !pattern.Match([]byte("2019-10-11T12:13:14-00:00")) {
		t.Fatal("Should match the time string.")
	}

	if !pattern.Match([]byte("2019-10-11T12:13:14+12:34")) {
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

func asPtr(str string) *string {
	return &str
}
func TestParseTime_rfc3339(t *testing.T) {
	expected := time.Date(2019, 3, 17, 10, 11, 12, 0, time.UTC)
	actual, err := parseTime("2019-03-17T10:11:12+00:00")
	if err != nil {
		t.Fatal(err)
	}

	expectedMillis := float64(expected.UnixNano()) / float64(time.Millisecond)
	assertEqual(expectedMillis, actual, 0.001, t)
}

func TestParseTime_timestamp(t *testing.T) {
	expected := 1552785731123.112
	actual, err := parseTime("1552785731123.112")
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(expected, actual, 0.001, t)
}

func TestParseTime_manualString(t *testing.T) {
	testString := func(str string, nowOffset, tolerance float64) {
		expected := currentTimeMillis() + nowOffset
		actual, err := parseTime(str)
		if err != nil {
			t.Fatal("Failed to parse time string:", str, err.Error())
		}
		assertEqual(expected, actual, tolerance, t)
	}

	testString("+0s", 0.0, 100.0) // now, within 100ms
	testString("+0days", 0.0, 100.0)

	testString("-1s", -1000.0, 100.0)
	testString("+1second", 1000.0, 100.0)
	testString("+2seconds", 2000.0, 100.0)

	testString("+1m", 1000*60.0, 100.0)
	testString("+1min", 1000*60.0, 100.0)
	testString("-3minutes", -3*1000*60.0, 100.0)

	testString("-1h", -1000*60*60.0, 100.0)
	testString("+4hrs", 4*1000*60*60.0, 100.0)
	testString("+1hour", 1000*60*60.0, 100.0)
	testString("+3hours", 3*1000*60*60.0, 100.0)

	testString("+2d", 2*1000*60*60*24.0, 100.0)
	testString("-2days", -2*1000*60*60*24.0, 100.0)
	testString("+1day", 1000*60*60*24.0, 100.0)

	testString("-1w", -1000*60*60*24*7.0, 100.0)
	testString("-3weeks", -3*1000*60*60*24*7.0, 100.0)
	testString("+1week", 1000*60*60*24*7.0, 100.0)

	testString("+1mn", 1000*60*60*24*30.0, 100.0)
	testString("+2mn", 2*1000*60*60*24*30.0, 100.0)
	testString("-1month", -1000*60*60*24*30.0, 100.0)
	testString("-3months", -3*1000*60*60*24*30.0, 100.0)

	testString("-1y", -1000*60*60*24*365.0, 100.0)
	testString("+1y", 1000*60*60*24*365.0, 100.0)
	testString("-2yr", -2*1000*60*60*24*365.0, 100.0)
	testString("-3years", -3*1000*60*60*24*365.0, 100.0)
}
