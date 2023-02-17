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
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/meblum/plane/cmd"
)

func main() {

	greetCommandName := "greet"
	greetType := "hello"
	greetName := "Anonymous"

	versionCommandName := "version"
	versionVerbose := false

	greetFlagSet := flag.NewFlagSet(greetCommandName, flag.ExitOnError)
	greetFlagSet.StringVar(&greetType, "type", greetType, "greet type (hello|bye)")
	greetFlagSet.StringVar(&greetName, "name", greetName, "name to print greeting for")

	versionFlagSet := flag.NewFlagSet(versionCommandName, flag.ExitOnError)
	versionFlagSet.BoolVar(&versionVerbose, "verbose", versionVerbose, "output version with additional information")

	c := &cmd.CommandSet{}
	c.Add("print greeting", greetFlagSet, false)
	c.Add("print version", versionFlagSet, false)

	currSub, _ := c.Parse(os.Args[1:], flag.ExitOnError)

	switch currSub.FlagSet.Name() {
	case greetCommandName:
		handleGreet(greetType, greetName)
	case versionCommandName:
		handleVersion("1.0.0", versionVerbose)
	}
}

func handleGreet(greetType, name string) {
	switch greetType {
	case "hello":
		fmt.Printf("Hello, %v! Nice to meet you :)", name)
	case "bye":
		fmt.Printf("Hello, %v! Nice to meet you :)", name)
	default:
		fmt.Print("invalid greet type, please use 'hello' or 'bye'")
	}
}

func handleVersion(version string, verbose bool) {
	if verbose {
		fmt.Printf("You are using Greeter version %v", version)
		return
	}
	fmt.Print(version)
}

```

See the [documentation][godev] for more information.



[godev]: https://pkg.go.dev/github.com/meblum/cmd



## License

MIT - See [LICENSE][license] file

[license]: https://github.com/meblum/cmd/blob/master/LICENSE
