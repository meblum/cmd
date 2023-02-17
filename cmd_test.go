package cmd

import (
	"flag"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCommandSet_Add(t *testing.T) {
	cmd := &CmdSet{}
	names := []string{"name", "andAnotherOne", "anotherName"}
	for _, v := range names {
		if c := cmd.Add("", flag.NewFlagSet(v, flag.ContinueOnError), false); cmd.commands[v] != c {
			t.Errorf("%v has not been added to commandSet", v)
		}
	}

	expectedLen := len(names[1])
	if cmd.cmdNameLength != expectedLen {
		t.Errorf("expected cmdNameLength length to be %v, got %v", expectedLen, cmd.cmdNameLength)
	}

	if len(cmd.commands) != 3 {
		t.Errorf("expected commands length to be %v, got %v", 3, cmd.cmdNameLength)
	}

}

func TestCommandSet_AddEmptyName(t *testing.T) {
	cmd := &CmdSet{}
	defer func() { recover() }()
	cmd.Add("", flag.NewFlagSet("", flag.ContinueOnError), false)
	t.Errorf("expected empty name to panic")
}

func TestCommandSet_AddExistingName(t *testing.T) {
	cmd := &CmdSet{}
	defer func() { recover() }()
	cmd.Add("name", flag.NewFlagSet("", flag.ContinueOnError), false)
	cmd.Add("name", flag.NewFlagSet("", flag.ContinueOnError), false)
	t.Errorf("expected empty name to panic")
}

func TestCommandSet_Visit(t *testing.T) {
	cmd := &CmdSet{}
	sets := map[*flag.FlagSet]bool{
		flag.NewFlagSet("a", flag.ContinueOnError): false,
		flag.NewFlagSet("b", flag.ContinueOnError): false,
		flag.NewFlagSet("c", flag.ContinueOnError): false,
	}
	for k := range sets {
		cmd.Add("", k, false)
	}

	cmd.Visit(func(s *Cmd) { sets[s.FlagSet] = true })
	for k, v := range sets {
		if !v {
			t.Errorf("%v has not been visited", k.Name())
		}
	}

}

func TestCommandSet_PrintUsage(t *testing.T) {

	emptyCmd := &CmdSet{output: &strings.Builder{}}

	emptyUsageCmd := &CmdSet{output: &strings.Builder{}}
	emptyUsageCmd.Add("", flag.NewFlagSet("a", flag.ContinueOnError), false)

	regularCmd := &CmdSet{output: &strings.Builder{}}
	regularCmd.Add("does b", flag.NewFlagSet("b", flag.ContinueOnError), false)
	regularCmd.Add("does a", flag.NewFlagSet("a", flag.ContinueOnError), false)

	tests := []struct {
		name     string
		c        *CmdSet
		expected string
	}{
		{
			"empty set",
			emptyCmd,
			"available subcommands for cmd.test.exe:\nuse \"<subcommand> --help\" for availble options of the specififc command",
		},
		{
			"empty usage",
			emptyUsageCmd,
			"available subcommands for cmd.test.exe:\n\ta - \nuse \"<subcommand> --help\" for availble options of the specififc command",
		},
		{
			"regular cmd",
			regularCmd,
			"available subcommands for cmd.test.exe:\n\ta - does a\n\tb - does b\nuse \"<subcommand> --help\" for availble options of the specififc command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.PrintUsage()
			if r := tt.c.output.(*strings.Builder).String(); r != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, r)
			}
		})
	}
}

func TestCommandSet_Parse(t *testing.T) {
	cmd := &CmdSet{}

	var av string
	af := flag.NewFlagSet("a", flag.ContinueOnError)
	af.StringVar(&av, "av", "", "")
	a := cmd.Add("", af, true)

	if c, err := cmd.Parse([]string{"a", "-av=aVal"}, flag.ContinueOnError); c != a || err != nil {
		t.Errorf("expected %v, got %v", a, c)
	}
	if av != "aVal" {
		t.Errorf("expected aVal, got %v", av)
	}

}

func TestCommandSet_ParseError(t *testing.T) {
	cmd := &CmdSet{}
	cmd.Add("", flag.NewFlagSet("a", flag.ContinueOnError), false)

	if os.Getenv("EXIT") == "EMPTY" {
		cmd.Parse([]string{}, flag.ExitOnError)
		return
	}

	if os.Getenv("EXIT") == "NO_MATCH" {
		cmd.Parse([]string{"foo"}, flag.ExitOnError)
		return
	}
	if os.Getenv("EXIT") == "EXTRA_ARGS" {
		cmd.Parse([]string{"a", "b"}, flag.ExitOnError)
		return
	}

	tests := []struct {
		name     string
		args     []string
		execFlag string
	}{
		{"no args", []string{}, "EMPTY"},
		{"no match", []string{"foo"}, "NO_MATCH"},
		{"extra args", []string{"a", "b"}, "EXTRA_ARGS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() { recover() }()
				_, err := cmd.Parse(tt.args, flag.PanicOnError)
				t.Errorf("expected panic, got %v", err)
			}()

			c := exec.Command(os.Args[0], "-test.run=TestCommandSet_ParseError")
			c.Env = append(os.Environ(), "EXIT="+tt.execFlag)
			err := c.Run()
			if e, ok := err.(*exec.ExitError); !ok || e.ExitCode() != 2 {
				t.Errorf("process ran with err %v, want exit status 1", err)
			}

			if _, err := cmd.Parse(tt.args, flag.ContinueOnError); err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestCommandSet_ParseHelp(t *testing.T) {
	cmd := &CmdSet{}

	if os.Getenv("EXIT") == "HELP" {
		cmd.Parse([]string{"-HeLP"}, flag.ExitOnError)
		return
	}

	c := exec.Command(os.Args[0], "-test.run=TestCommandSet_ParseHelp")
	c.Env = append(os.Environ(), "EXIT=HELP")

	if err := c.Run(); err != nil {
		t.Errorf("process ran with err %v, want exit status 0", err)
	}

	if _, err := cmd.Parse([]string{"-HeLP"}, flag.ContinueOnError); err != flag.ErrHelp {
		t.Errorf("process ran with err %v, want flag.ErrHelp", err)
	}

}
