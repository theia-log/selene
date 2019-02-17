package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Command func(args []string) error

type commandData struct {
	help    string
	command Command
}

type CommandsCLI struct {
	prog        string
	description string
	commands    map[string]*commandData
	out         io.Writer
	err         io.Writer
}

func NewCommandCLI(prog, description string) *CommandsCLI {
	return &CommandsCLI{
		prog:        prog,
		description: description,
		commands:    map[string]*commandData{},
		out:         os.Stdout,
		err:         os.Stderr,
	}
}

func (c *CommandsCLI) AddCommand(name string, command Command, help string) *CommandsCLI {
	c.commands[name] = &commandData{
		help:    help,
		command: command,
	}
	return c
}

func (c *CommandsCLI) Execute() error {
	return c.ExecuteWithArgs(os.Args[1:])
}

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
	return commands
}

func (c *CommandsCLI) printHelpTo(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s %s [ARGS]\n\n%s\n\n", c.prog, strings.Join(c.getCommandNames(), "|"), c.description)
	fmt.Fprintf(out, "Commands:\n")
	for command, descriptor := range c.commands {
		fmt.Fprintf(out, "\t%s - %s\n", command, descriptor.help)
	}
}

func (c *CommandsCLI) PrintHelp() {
	c.printHelpTo(c.err)
}

func (c *CommandsCLI) UseOut(out io.Writer) *CommandsCLI {
	c.out = out
	return c
}

func (c *CommandsCLI) UseErr(err io.Writer) *CommandsCLI {
	c.err = err
	return c
}
