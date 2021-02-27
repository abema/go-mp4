package main

import (
	"fmt"
	"os"

	"github.com/abema/go-mp4/mp4tool/divide"
	"github.com/abema/go-mp4/mp4tool/dump"
	"github.com/abema/go-mp4/mp4tool/edit"
	"github.com/abema/go-mp4/mp4tool/probe"
	"github.com/abema/go-mp4/mp4tool/psshdump"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "help":
		printUsage()
	case "dump":
		dump.Main(args[1:])
	case "psshdump":
		psshdump.Main(args[1:])
	case "probe":
		probe.Main(args[1:])
	case "alpha":
		alpha(args[1:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func alpha(args []string) {
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "edit":
		edit.Main(args[1:])
	case "divide":
		divide.Main(args[1:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "USAGE: mp4tool COMMAND_NAME [ARGS]\n")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "COMMAND_NAME:")
	fmt.Fprintln(os.Stderr, "  dump")
	fmt.Fprintln(os.Stderr, "  psshdump")
	fmt.Fprintln(os.Stderr, "  probe")
	fmt.Fprintln(os.Stderr, "  alpha edit")
	fmt.Fprintln(os.Stderr, "  alpha divide")
}
