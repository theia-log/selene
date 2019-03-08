package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewCommandCLI(t *testing.T) {
	cli := NewCommandCLI("test-prog", "Test Program")

	var out bytes.Buffer
	cli.UseErr(&out)

	cli.PrintHelp()

	const expected string = "Usage: test-prog  [ARGS]\n\nTest Program\n"

	if out.String() != expected {
		t.Fatalf("Expected '%s' but got '%s'\n", expected, out.String())
	}
}

func TestNewCommandCLI_withCommands(t *testing.T) {
	executedCommands := map[string]string{}
	cli := NewCommandCLI("test-prog", "Test Program")

	cli.AddCommand("command1", func(args []string) error {
		executedCommands["command1"] = "yes"
		if args == nil || len(args) != 1 {
			t.Fatal("Expected to receive 1 argument for command 1")
		}
		if args[0] != "arg1" {
			t.Fatal("Got invalid value for argument 1")
		}
		return nil
	}, "Command 1")

	cli.AddCommand("command2", func(args []string) error {
		executedCommands["command2"] = "yes"
		if args == nil || len(args) != 1 {
			t.Fatal("Expected to receive 1 argument for command 1")
		}
		if args[0] != "arg2" {
			t.Fatal("Got invalid value for argument 1")
		}
		return nil
	}, "Command 2")

	var out bytes.Buffer
	cli.UseErr(&out)

	cli.PrintHelp()

	expected := strings.Join([]string{
		"Usage: test-prog command1|command2 [ARGS]\n",
		"Test Program\n",
		"Commands:",
		"\tcommand1 - Command 1",
		"\tcommand2 - Command 2",
		"",
	}, "\n")

	if out.String() != expected {
		t.Fatalf("Expected '%s' but got '%s'\n", expected, out.String())
	}

	cli.ExecuteWithArgs([]string{"command1", "arg1"})

	if _, ok := executedCommands["command1"]; !ok {
		t.Fatal("Command 1 not exected, but it should have been")
	}

	cli.ExecuteWithArgs([]string{"command2", "arg2"})

	if _, ok := executedCommands["command2"]; !ok {
		t.Fatal("Command 2 not exected, but it should have been")
	}
}
