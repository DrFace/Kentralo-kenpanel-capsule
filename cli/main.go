package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kentralo/kenpanel-capsule/packager"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "inspect":
		inspectCmd := flag.NewFlagSet("inspect", flag.ExitOnError)
		dir := inspectCmd.String("dir", ".", "Application directory to inspect")
		inspectCmd.Parse(os.Args[2:])

		p := packager.NewPackager(*dir)
		m, err := p.Inspect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error inspecting application: %v\n", err)
			os.Exit(1)
		}

		data, _ := json.MarshalIndent(m, "", "  ")
		fmt.Println(string(data))

	case "version":
		fmt.Println("kenpanel-capsule v2.3.0 (standalone engine)")

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: kenpanel-capsule <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  inspect     Analyze application runtime, ports, dependencies and redact secrets")
	fmt.Println("  pack        Bundle application into portable .tar.zst archive with manifest")
	fmt.Println("  version     Show version")
}
