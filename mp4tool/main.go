package main

import (
	"fmt"
	"os"

	"github.com/abema/go-mp4/mp4tool/divide"
	"github.com/abema/go-mp4/mp4tool/dump"
	"github.com/abema/go-mp4/mp4tool/edit"
	"github.com/abema/go-mp4/mp4tool/extract"
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
		os.Exit(dump.Main(args[1:]))
	case "psshdump":
		os.Exit(psshdump.Main(args[1:]))
	case "probe":
		os.Exit(probe.Main(args[1:]))
	case "extract":
		os.Exit(extract.Main(args[1:]))
	case "alpha":
		os.Exit(alpha(args[1:]))
	default:
		printUsage()
		os.Exit(1)
	}
}

func alpha(args []string) int {
	if len(args) < 1 {
		printUsage()
		return 1
	}

	switch args[0] {
	case "edit":
		return edit.Main(args[1:])
	case "divide":
		return divide.Main(args[1:])
	default:
		printUsage()
		return 1
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "USAGE: mp4tool COMMAND_NAME [ARGS]\n")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "COMMAND_NAME:")
	fmt.Fprintln(os.Stderr, "  dump         : display box tree as human readable format")
	fmt.Fprintln(os.Stderr, "  psshdump     : display pssh box attributes")
	fmt.Fprintln(os.Stderr, "  probe        : probe and summarize mp4 file status")
	fmt.Fprintln(os.Stderr, "  extract      : extract specific box")
	fmt.Fprintln(os.Stderr, "  alpha edit")
	fmt.Fprintln(os.Stderr, "  alpha divide")
}
