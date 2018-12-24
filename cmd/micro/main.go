package main

import (
	"fmt"
	"os"

	"github.com/micro/enterprise/go/license"
	"github.com/micro/enterprise/go/plugin"
	"github.com/micro/enterprise/go/token"
	"github.com/micro/go-micro/cmd"
	dmc "github.com/micro/micro/cmd"
	mp "github.com/micro/micro/plugin"
)

var (
	name        = "micro"
	description = "An enterprise microservice toolkit"
	version     = "0.1.0"
)

func main() {
	// register plugin
	if err := mp.Register(license.NewPlugin()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := mp.Register(plugin.NewPlugin()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// setup the command line
	dmc.Setup(cmd.App())

	// add commands
	app := cmd.App()
	app.Commands = append(app.Commands, token.Commands()...)
	app.Commands = append(app.Commands, license.Commands()...)
	app.Commands = append(app.Commands, plugin.Commands()...)

	// initialise command line
	cmd.Init(
		cmd.Name(name),
		cmd.Description(description),
		cmd.Version(version),
	)
}
