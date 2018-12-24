package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	"github.com/micro/micro/plugin"
)

func buildSo(soPath string, parts []string) error {
	// check if .so file exists
	if _, err := os.Stat(soPath); os.IsExist(err) {
		return nil
	}

	// name and things
	name := parts[len(parts)-1]
	// type of plugin
	typ := parts[0]
	// new func signature
	newfn := fmt.Sprintf("New%s", strings.Title(typ))

	// micro has NewPlugin type def
	if typ == "micro" {
		newfn = "NewPlugin"
	}

	// now build the plugin
	if err := Build(soPath, &Plugin{
		Name:    name,
		Type:    typ,
		Path:    filepath.Join(append([]string{"github.com/micro/go-plugins"}, parts...)...),
		NewFunc: newfn,
	}); err != nil {
		return fmt.Errorf("Failed to build plugin %s: %v", name, err)
	}

	return nil
}

func load(p string) error {
	p = strings.TrimSpace(p)

	if len(p) == 0 {
		return nil
	}

	parts := strings.Split(p, "/")

	// 1 part means local plugin
	// plugin/foobar
	if len(parts) == 1 {
		return fmt.Errorf("Unknown plugin %s", p)
	}

	// set soPath to specified path
	soPath := p

	// build on the fly if not .so
	if !strings.HasSuffix(p, ".so") {
		// set new so path
		soPath = filepath.Join("plugin", p+".so")

		// build new .so
		if err := buildSo(soPath, parts); err != nil {
			return err
		}
	}

	// load the plugin
	pl, err := Load(soPath)
	if err != nil {
		return fmt.Errorf("Failed to load plugin %s: %v", soPath, err)
	}

	switch pl.Type {
	case "micro":
		pg, ok := pl.NewFunc.(func() plugin.Plugin)
		if !ok {
			return fmt.Errorf("Invalid plugin %s", pl.Name)
		}
		plugin.Register(pg())
	case "broker":
		pg, ok := pl.NewFunc.(func(...broker.Option) broker.Broker)
		if !ok {
			return fmt.Errorf("Invalid plugin %s", pl.Name)
		}
		cmd.DefaultBrokers[pl.Name] = pg
	case "client":
		pg, ok := pl.NewFunc.(func(...client.Option) client.Client)
		if !ok {
			return fmt.Errorf("Invalid plugin %s", pl.Name)
		}
		cmd.DefaultClients[pl.Name] = pg
	case "registry":
		pg, ok := pl.NewFunc.(func(...registry.Option) registry.Registry)
		if !ok {
			return fmt.Errorf("Invalid plugin %s", pl.Name)
		}
		cmd.DefaultRegistries[pl.Name] = pg

	case "selector":
		pg, ok := pl.NewFunc.(func(...selector.Option) selector.Selector)
		if !ok {
			return fmt.Errorf("Invalid plugin %s", pl.Name)
		}
		cmd.DefaultSelectors[pl.Name] = pg
	case "server":
		pg, ok := pl.NewFunc.(func(...server.Option) server.Server)
		if !ok {
			return fmt.Errorf("Invalid plugin %s", pl.Name)
		}
		cmd.DefaultServers[pl.Name] = pg
	case "transport":
		pg, ok := pl.NewFunc.(func(...transport.Option) transport.Transport)
		if !ok {
			return fmt.Errorf("Invalid plugin %s", pl.Name)
		}

		cmd.DefaultTransports[pl.Name] = pg
	default:
		return fmt.Errorf("Unknown plugin type: %s for %s", pl.Type, pl.Name)
	}

	return nil
}

// returns a micro plugin which loads plugins
func NewPlugin() plugin.Plugin {
	return plugin.NewPlugin(
		plugin.WithName("plugins"),
		plugin.WithFlag(
			cli.StringSliceFlag{
				Name:   "plugins",
				EnvVar: "MICRO_PLUGINS",
				Usage:  "Comma separated list of plugins e.g broker/rabbitmq, registry/etcd, micro/basic_auth, /path/to/plugin.so",
			},
		),
		plugin.WithInit(func(ctx *cli.Context) error {
			plugins := ctx.StringSlice("plugins")
			if len(plugins) == 0 {
				return nil
			}

			for _, p := range plugins {
				if err := load(p); err != nil {
					return err
				}
				log.Logf("Loaded plugin %s\n", p)
			}

			return nil
		}),
	)
}
