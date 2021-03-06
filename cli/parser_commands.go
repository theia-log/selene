package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Command represents a command runner. It takes subset of the program arguments
// that apply to the parsed command.
type Command func(args []string) error

// commandData wraps the command handler and the help data. Used for holding the
// program help and is useful while starting the program and parsing the
// command line.
type commandData struct {
	help    string
	command Command
}

// CommandsCLI is an implementation of a program CLI based on sub-commands.
// Parses the command line arguments into subcommands and then executes the
// command given on the command line.
//
// This is a programmable API for writing programs that use subcommands.
//
type CommandsCLI struct {
	prog        string
	description string
	commands    map[string]*commandData
	out         io.Writer
	err         io.Writer
}

// NewCommandCLI builds a new CommandsCLI for the given program name and
// description.
func NewCommandCLI(prog, description string) *CommandsCLI {
	return &CommandsCLI{
		prog:        prog,
		description: description,
		commands:    map[string]*commandData{},
		out:         os.Stdout,
		err:         os.Stderr,
	}
}

// AddCommand add (sub)command.
// A command is given a name, a handler (Command handler) and a help string.
func (c *CommandsCLI) AddCommand(name string, command Command, help string) *CommandsCLI {
	c.commands[name] = &commandData{
		help:    help,
		command: command,
	}
	return c
}

// Execute executes the CLI  program with the default arguments. The arguments
// are obtained from the program invocation (passed by the OS). The zeroth
// argument is ignored (as this is the program name), and the first argument
// is used to determine the command name.
func (c *CommandsCLI) Execute() error {
	return c.ExecuteWithArgs(os.Args[1:])
}

// ExecuteWithArgs executes the CLI  program with the given arguments.
// The first argument must be the command name. The rest of the arguments are
// going to be passed down to the appropriate command handler.
// If no hander exists for the command, an error is returned.
func (c *CommandsCLI) ExecuteWithArgs(args []string) error {

	if args == nil || len(args) == 0 {
		c.PrintHelp()
		return nil
	}
	commandName := args[0]
	remArgs := args[1:]

	descriptor, ok := c.commands[commandName]
	if !ok {
		fmt.Fprintf(c.err, "%s: unknown command '%s'\n", c.prog, commandName)
		c.PrintHelp()
		return nil
	}
	return descriptor.command(remArgs)
}

func (c *CommandsCLI) getCommandNames() []string {
	commands := []string{}
	for command := range c.commands {
		commands = append(commands, command)
	}
	sort.Slice(commands, func(i, j int) bool { return commands[i] < commands[j] })
	return commands
}

func (c *CommandsCLI) printHelpTo(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s %s [ARGS]\n\n%s\n", c.prog, strings.Join(c.getCommandNames(), "|"), c.description)
	if c.commands != nil && len(c.commands) > 0 {
		fmt.Fprintf(out, "\nCommands:\n")
		for _, command := range c.getCommandNames() {
			descriptor := c.commands[command]
			fmt.Fprintf(out, "\t%s - %s\n", command, descriptor.help)
		}
	}
}

// PrintHelp prints the help string to the set output stream.
func (c *CommandsCLI) PrintHelp() {
	c.printHelpTo(c.err)
}

// UseOut replaces the output writer used (which by default is os.Stdout) with
// the output writer provided as argument. All messages generated by the
// command CLI (like printing program help) are going to be written to this
// output writer.
func (c *CommandsCLI) UseOut(out io.Writer) *CommandsCLI {
	c.out = out
	return c
}

// UseErr replaces the error output writer used (which by default is os.Stderr)
// with the provided writer. All error messages are going to be written to this
// writer (like standard error during parsing).
func (c *CommandsCLI) UseErr(err io.Writer) *CommandsCLI {
	c.err = err
	return c
}
