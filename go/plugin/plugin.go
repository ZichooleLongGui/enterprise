package plugin

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"
	"text/template"
)

// Plugin is a plugin loaded from a file
type Plugin struct {
	// Name of the plugin e.g rabbitmq
	Name string
	// Type of the plugin e.g broker
	Type string
	// Path specifies the import path
	Path string
	// NewFunc creates an instance of the plugin
	NewFunc interface{}
}

// Load loads a plugin created with `go build -buildmode=plugin`
func Load(path string) (*Plugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}
	pl, ok := s.(*Plugin)
	if !ok {
		return nil, errors.New("could not find plugin")
	}
	return pl, nil
}

// Generate creates a go file at the specified path.
// You must use `go build -buildmode=plugin`to build it.
func Generate(path string, p *Plugin) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	t, err := template.New(p.Name).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(f, p)
}

// Build generates a dso plugin using the go command `go build -buildmode=plugin`
func Build(path string, p *Plugin) error {
	path = strings.TrimSuffix(path, ".so")

	// create go file in tmp path
	temp := os.TempDir()
	base := filepath.Base(path)
	goFile := filepath.Join(temp, base+".go")

	// generate .go file
	if err := Generate(goFile, p); err != nil {
		return err
	}
	// remove .go file
	defer os.Remove(goFile)

	c := exec.Command("go", "build", "-buildmode=plugin", "-o", path+".so", goFile)
	return c.Run()
}
