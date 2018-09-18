package main

import (
	"github.com/micro/enterprise/go/license"
	"github.com/micro/enterprise/go/token"
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

	// add commands
	app := cmd.App()
	app.Commands = append(app.Commands, token.Commands()...)
	app.Commands = append(app.Commands, license.Commands()...)

	// initialise command line
	cmd.Init(
		cmd.Name(name),
		cmd.Description(description),
		cmd.Version(version),
	)
}
