package main

import (
	"experts-meeting/cli"
	"fmt"
	"os"
)

func main() {
	if err := cli.CLIRun(os.Args); err != nil {
		fmt.Println("Error:", err)
		// os.Exit(1)
	}
}
