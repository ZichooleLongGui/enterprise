// Package config provides dynamic config via env vars, flags and config file
package config

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-config/source/microcli"
)

// NewConfig returns new config for env vars, flags and config file.
// Config file is expected to be at the path ./config.json.
func NewConfig(opts ...config.Option) config.Config {
	return config.NewConfig(
		// base config from env
		config.WithSource(env.NewSource()),
		// override env with flags
		config.WithSource(microcli.NewSource()),
		// override flags with file
		config.WithSource(file.NewSource()),
	)
}
