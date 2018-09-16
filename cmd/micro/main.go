package main

import (
	"github.com/micro/go-micro/cmd"
	dmc "github.com/micro/micro/cmd"
)

var (
	name        = "micro"
	description = "An enterprise cloud-native toolkit"
	version     = "0.1.0"
)

func main() {
	// setup the command line
	dmc.Setup(cmd.App())

	// initialise command line
	cmd.Init(
		cmd.Name(name),
		cmd.Description(description),
		cmd.Version(version),
	)
}
