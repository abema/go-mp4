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
		return
	}

	switch args[0] {
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
	}
}

func alpha(args []string) {
	if len(args) < 1 {
		printUsage()
		return
	}

	switch args[0] {
	case "edit":
		edit.Main(args[1:])
	case "divide":
		divide.Main(args[1:])
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Printf("USAGE: mp4tool COMMAND_NAME [ARGS]\n")
	fmt.Println()
	fmt.Println("COMMAND_NAME:")
	fmt.Println("  dump")
	fmt.Println("  psshdump")
	fmt.Println("  probe")
	fmt.Println("  alpha edit")
	fmt.Println("  alpha divide")
}
