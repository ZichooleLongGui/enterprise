package main

import (
	"fmt"
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
	dmc "github.com/micro/micro/cmd"
	mp "github.com/micro/micro/plugin"

	// enterprise plugins
	"github.com/micro/enterprise/go/auth"
	"github.com/micro/enterprise/go/license"
	"github.com/micro/enterprise/go/metrics"
	"github.com/micro/enterprise/go/plugin"
	"github.com/micro/enterprise/go/token"

	// TODO: move to plugin dir
	_ "github.com/micro/go-plugins/micro/bot/input/discord"
	_ "github.com/micro/go-plugins/micro/bot/input/telegram"
	"github.com/micro/go-plugins/micro/cors"

	// grpc by default
	gbkr "github.com/micro/go-plugins/broker/grpc"
	gcli "github.com/micro/go-plugins/client/grpc"
	gsrv "github.com/micro/go-plugins/server/grpc"
)

var (
	name        = "micro"
	description = "An enterprise microservice toolkit"
	version     = "0.2.0"
)

// TODO: move to plugin/ dir
func plugins() {
	// register license plugin
	if err := mp.Register(license.NewPlugin()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// register plugin loader
	if err := mp.Register(plugin.NewPlugin()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// register admin auth
	if err := mp.Register(auth.NewPlugin()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// register cors plugin
	if err := mp.Register(cors.NewPlugin()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// register metrics
	if err := mp.Register(metrics.NewPlugin()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set grpc defaults
	broker.DefaultBroker = gbkr.NewBroker()
	client.DefaultClient = gcli.NewClient()
	server.DefaultServer = gsrv.NewServer()
}

func commands(app *cli.App) {
	// add commands
	app.Commands = append(app.Commands, token.Commands()...)
	app.Commands = append(app.Commands, license.Commands()...)
	app.Commands = append(app.Commands, plugin.Commands()...)
}

func main() {
	// setup plugins
	plugins()

	// setup the command line
	dmc.Setup(cmd.App())

	// setup plugins
	commands(cmd.App())

	// initialise command line
	cmd.Init(
		cmd.Name(name),
		cmd.Description(description),
		cmd.Version(version),
	)
}
