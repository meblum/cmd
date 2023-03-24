package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Handler interface {
	HandleCmd(*Cmd) error
}

// A Cmd is a single subcommand with a set of flags.
type Cmd struct {
	Info      string
	FlagSet   *flag.FlagSet
	AllowArgs bool
	Handler   Handler
}

func (c *Cmd) Handle() error {
	return c.Handler.HandleCmd(c)
}

// A CmdSet contains a set of subcommands.
//
// The zero value of CmdSet contains no commands.
// Use CmdSet.Add to add subcommands.
type CmdSet struct {
	commands map[string]*Cmd
	// usage print destination. Defaults to os.Stderr
	output io.Writer
	// used to pad usage info
	cmdNameLength int
}

// Add adds a subcommand with specified usage string and flag set.
// The command name is derived from flags.Name.
// Use the allowArgs argument to specify wheather additional args should be allowed.
// For easy mapping from command to the relevant handler, supply handler which will be added to the Cmd returned from Parse().
// Returns the addded command.
//
// Command names must be unique within a CommandSet.
// An attempt to define a command whose name is an emty string or whose name is already in use will cause a panic.
func (c *CmdSet) Add(usage string, flags *flag.FlagSet, handler Handler, allowArgs bool) *Cmd {
	if c.commands == nil {
		c.commands = make(map[string]*Cmd)
	}

	cmdName := flags.Name()
	if _, ok := c.commands[cmdName]; ok || cmdName == "" {
		panic("invalid command name " + cmdName)
	}

	cmdNameLen := len(flags.Name())
	if cmdNameLen > c.cmdNameLength {
		c.cmdNameLength = cmdNameLen
	}
	cmd := &Cmd{Info: usage, AllowArgs: allowArgs, FlagSet: flags}
	c.commands[cmdName] = cmd
	return cmd
}

// Visit runs f on every command currently in c.
// Most commonly used to apply global flags to a command set.
func (c *CmdSet) Visit(f func(*Cmd)) {
	for _, v := range c.commands {
		f(v)
	}
}

// PrintUsage prints usage information to standard error.
func (c *CmdSet) PrintUsage() {
	output := c.output
	if output == nil {
		output = os.Stderr
	}
	cliName := filepath.Base(os.Args[0])
	fmt.Fprintf(output, "available subcommands for %v:\n", cliName)
	padVerb := fmt.Sprintf("%%-%vs", c.cmdNameLength)

	for _, v := range c.sortedCommandNames() {
		paddedCmdName := fmt.Sprintf(padVerb, v)
		fmt.Fprintf(output, "\t%v - %v\n", paddedCmdName, c.commands[v].Info)
	}

	fmt.Fprint(output, "use \"<subcommand> --help\" for availble options of the specififc command")
}

func (c *CmdSet) sortedCommandNames() []string {
	var keys []string
	for k := range c.commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (c *CmdSet) getSubcommand(arguments []string) (*Cmd, error) {
	if len(arguments) < 1 {
		return nil, fmt.Errorf("subcommand not specified")
	}
	requestedSubcommand := arguments[0]
	trimmed := strings.TrimLeft(requestedSubcommand, "-")
	if strings.EqualFold(trimmed, "h") || strings.EqualFold(trimmed, "help") {
		return nil, flag.ErrHelp
	}

	var subcommand *Cmd
	for k, v := range c.commands {
		if strings.EqualFold(requestedSubcommand, k) {
			subcommand = v
			break
		}
	}

	if subcommand == nil {
		return nil, fmt.Errorf("invalid subcommand %q", requestedSubcommand)
	}

	return subcommand, nil
}

// Parse parses the subcommand from arguments[0] and its flags from arguments[1:] with the error handling specified by errorHandling.
// Returns the supplied Subcommand if a match was found.
//
// Must be called after all subcommands are defined and before flags are accessed by the program.
// If arguments is nil, will default to `os.Args[1:]`.
func (c *CmdSet) Parse(arguments []string, errorHandling flag.ErrorHandling) (*Cmd, error) {
	if arguments == nil {
		arguments = os.Args[1:]
	}

	subcommand, err := c.getSubcommand(arguments)

	if err != nil {
		c.PrintUsage()
		handleError(err, errorHandling)
		return nil, err
	}

	if err := subcommand.FlagSet.Parse(arguments[1:]); err != nil {
		return subcommand, err
	}

	if !subcommand.AllowArgs && subcommand.FlagSet.NArg() > 0 {
		subcommand.FlagSet.Usage()
		err := fmt.Errorf("arguments not supported - %v", subcommand.FlagSet.Args())
		handleError(err, errorHandling)
		return subcommand, err
	}

	return subcommand, nil
}

func handleError(err error, errorHandling flag.ErrorHandling) {
	switch errorHandling {
	case flag.ContinueOnError:
		return
	case flag.ExitOnError:
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(2)
	case flag.PanicOnError:
		panic(err)
	}
}
