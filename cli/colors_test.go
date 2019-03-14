package cli

import "testing"

func TestContextGetString(t *testing.T) {
	ctx := Context{
		"hasIt": "yes",
	}

	val := ctx.GetString("hasIt")
	if val != "yes" {
		t.Fatal("Expected to get the value from context.")
	}

	val = ctx.GetString("doesNotHaveIt")
	if val != "" {
		t.Fatal("Expected to get an empty string.")
	}
}
