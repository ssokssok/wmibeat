package main

import (
	"os"

	"github.com/ssokssok/wmibeat/cmd"

//	_ "github.com/ssokssok/wmibeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
