package license

import (
	"fmt"
	"os"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"
)

// returns a micro plugin which validates use of license
func NewPlugin() plugin.Plugin {
	return plugin.NewPlugin(
		plugin.WithName("license"),
		plugin.WithInit(func(ctx *cli.Context) error {
			if len(ctx.Args()) == 0 {
				return nil
			}

			name := ctx.Args()[0]

			// skip on enterprise commands
			if name == "license" || name == "token" {
				return nil
			}

			key := os.Getenv("MICRO_LICENSE_KEY")

			// TODO: check key is valid
			if len(key) < 62 {
				fmt.Println("Micro Enterprise license key missing")
				os.Exit(1)
			}

			return nil
		}),
	)
}
