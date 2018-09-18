package license

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/micro/cli"
)

func generate(ctx *cli.Context) {
	sub := ctx.String("subscription")

	if len(sub) == 0 {
		fmt.Println("Subscription is blank (specify --subscription)")
		os.Exit(1)
	}

	// generate
	l, err := Generate(sub)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Your license (set as MICRO_ENTERPRISE_LICENSE env var or X-Micro-License http header):")
	fmt.Println(l)
}

func revoke(ctx *cli.Context) {
	license := ctx.String("license")
	if len(license) == 0 {
		fmt.Println("Licence is blank (specify --license)")
		os.Exit(1)
	}

	// revoke license
	if err := Revoke(license); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Licence revoked")
}

func verify(ctx *cli.Context) {
	license := ctx.String("license")
	if len(license) == 0 {
		fmt.Println("Licence is blank (specify --license)")
		os.Exit(1)
	}

	// revoke license
	if err := Verify(license); err != nil {
		fmt.Println("Verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("Licence verified")
}

func list(ctx *cli.Context) {
	licenses, err := List()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(licenses) == 0 {
		fmt.Println(`{}`)
		return
	}
	j, err := json.MarshalIndent(licenses, "", "\t")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(j))
}

func subscriptions(ctx *cli.Context) {
	subs, err := Subscriptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(subs) == 0 {
		fmt.Println(`{}`)
		return
	}
	j, err := json.MarshalIndent(subs, "", "\t")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(j))
}

func licenseCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "list",
			Usage:  "List licenses",
			Action: list,
		},
		{
			Name:   "subscriptions",
			Usage:  "List subscriptions",
			Action: subscriptions,
		},
		{
			Name:   "generate",
			Usage:  "Generate an api license (specify --subscription)",
			Action: generate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "subscription",
					Usage: "Subcription to generate the license for",
				},
			},
		},
		{
			Name:   "revoke",
			Usage:  "Revoke an api license (specify --license)",
			Action: revoke,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "license",
					Usage: "Encoded license key to revoke",
				},
			},
		},
		{
			Name:   "verify",
			Usage:  "Verify an api license is valid (specify --license)",
			Action: verify,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "license",
					Usage: "Encoded license key to verify",
				},
			},
		},
	}
}

// Commands returns license commands
func Commands() []cli.Command {
	return []cli.Command{{
		Name:        "license",
		Usage:       "Enterprise license commands",
		Subcommands: licenseCommands(),
	}}
}
