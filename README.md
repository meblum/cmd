# Idiomatic Go Package for cli subcommands

[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)][godev]

This package is intended to be used in conjunction with the Go `flag` package for maping subcommands to their respective flagset.

The primary features of `cmd` are:

- Use the standerd `flag` packege. `cmd` does not implement its own flag parsing. Instead it integrates nicely with the standerd `flag` package.

- Common API, style and behavior closely follows that of the standers `flag` packege.

- Minimal feature set, does not attempt to be a "one size fits all" solution.

- Optionaly enable/disable extra command argument support (besides flags) for subcommands.

## Install

```
go get -u github.com/meblum/cmd
```

## Example usage

```go
import (
	"flag"
	"fmt"

	"github.com/meblum/cmd"
)

func main() {
	greeter := greetHandler{
		greetType: "hello",
		greetName: "Anonymous",
	}

	greetFlagSet := flag.NewFlagSet("greet", flag.ExitOnError)
	greetFlagSet.StringVar(&greeter.greetType, "type", greeter.greetType, "greet type (hello|bye)")
	greetFlagSet.StringVar(&greeter.greetName, "name", greeter.greetName, "name to print greeting for")

	version := versionHandler{
		version: "1.0.0",
		verbose: false,
	}

	versionFlagSet := flag.NewFlagSet("version", flag.ExitOnError)
	versionFlagSet.BoolVar(&version.verbose, "verbose", version.verbose, "output version with additional information")

	c := &cmd.CmdSet{}
	c.Add("print greeting", greetFlagSet, greeter, false)
	c.Add("print version", versionFlagSet, version, false)
	c.HandleCmd(nil, flag.ExitOnError)
}

type greetHandler struct {
	greetType string
	greetName string
}

func (g greetHandler) Handle(c *cmd.Cmd) error {
	switch g.greetType {
	case "hello":
		fmt.Printf("Hello, %v! Nice to meet you :)", g.greetName)
	case "bye":
		fmt.Printf("Hello, %v! Nice to meet you :)", g.greetName)
	default:
		return fmt.Errorf("invalid greet type %v, valid options are 'hello' or 'bye'", g.greetType)
	}
	return nil
}

type versionHandler struct {
	version string
	verbose bool
}

func (g versionHandler) Handle(c *cmd.Cmd) error {
	if g.verbose {
		fmt.Printf("You are using Greeter version %v", g.version)
	}
	fmt.Print(g.version)
	return nil
}


```

See the [documentation][godev] for more information.



[godev]: https://pkg.go.dev/github.com/meblum/cmd



## License

MIT - See [LICENSE][license] file

[license]: https://github.com/meblum/cmd/blob/master/LICENSE
