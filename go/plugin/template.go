package plugin

var (
	tmpl = `
package main

import (
	"{{.Path}}"
	"github.com/micro/enterprise/go/plugin"
)

var Plugin = &plugin.Plugin{
	Name: "{{.Name}}",
	Type: "{{.Type}}",
	Path: "{{.Path}}",
	NewFunc: {{.Name}}.{{.NewFunc}},
}
`
)
