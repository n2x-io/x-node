package main

import (
	"fmt"
	// "log"

	"n2x.dev/x-lib/pkg/version"
	"n2x.dev/x-node/internal/app/node/cmd"
)

func main() {
	// if err := cmd.ConsoleInit(); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Printf("%s %s ", version.NODE_NAME, version.GetVersion())

	cmd.Execute()
}
