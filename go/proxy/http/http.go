// Package http provides a go-micro to http proxy
package http

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	emicro "github.com/micro/enterprise/go/micro"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
)

// Router will proxy rpc requests as http POST requests. It is a server.Router
type Router struct {
	// Converts RPC Foo.Bar to /foo/bar
	Resolver *Resolver
	// The http endpoint to call
	Endpoint string

	// first request
	first bool

	// rpc ep / http ep mapping
	eps map[string]string
}

// Resolver resolves rpc to http. It explicity maps Foo.Bar to /foo/bar
type Resolver struct{}

var (
	// The default endpoint
	Endpoint = "http://localhost:9090"
)

// Foo.Bar becomes /foo/bar
func (r *Resolver) Resolve(ep string) string {
	// replace . with /
	ep = strings.Replace(ep, ".", "/", -1)
	// lowercase the whole thing
	ep = strings.ToLower(ep)
	// prefix with "/"
	return filepath.Join("/", ep)
}

// set the nil things
func (p *Router) setup() {
	if p.Resolver == nil {
		p.Resolver = new(Resolver)
	}
	if p.Endpoint == "" {
		p.Endpoint = Endpoint
	}
	if p.eps == nil {
		p.eps = map[string]string{}
	}
}

// GetEndpoint returns the http endpoint for an rpc endpoint
// GetEndpoint("Foo.Bar") returns /foo/bar
func (p *Router) GetEndpoint(rpcEp string) (string, error) {
	p.setup()

	// get http endpoint
	ep, ok := p.eps[rpcEp]
	if ok {
		return ep, nil
	}

	// get default
	ep = p.Resolver.Resolve(rpcEp)

	// full path to call
	u, err := url.Parse(p.Endpoint)
	if err != nil {
		return "", err
	}

	// set path
	u.Path = ep

	// create ep
	return u.String(), nil
}

// RegisterEndpoint registers a http endpoint against an RPC endpoint
//	RegisterEndpoint("Foo.Bar", "/foo/bar")
//	RegisterEndpoint("Greeter.Hello", "/helloworld")
//	RegisterEndpoint("Greeter.Hello", "http://localhost:8080/")
func (p *Router) RegisterEndpoint(rpcEp string, httpEp string) error {
	p.setup()

	// register if already http
	if strings.HasPrefix(httpEp, "http://") || strings.HasPrefix(httpEp, "https://") {
		p.eps[rpcEp] = httpEp
		return nil
	}

	// full path to call
	u, err := url.Parse(p.Endpoint)
	if err != nil {
		return err
	}

	// set path
	u.Path = httpEp

	// create ep
	p.eps[rpcEp] = u.String()

	return nil
}

// ServeRequest honours the server.Router interface
func (p *Router) ServeRequest(ctx context.Context, req server.Request, rsp server.Response) error {
	// rudimentary post based streaming
	for {
		// get data
		body, err := req.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		var rpcEp string

		// get rpc endpoint
		if p.first {
			p.first = false
			rpcEp = req.Endpoint()
		} else {
			hdr := req.Header()
			rpcEp = hdr["X-Micro-Endpoint"]
		}

		// get http endpoint
		ep, err := p.GetEndpoint(rpcEp)
		if err != nil {
			return errors.NotFound(req.Service(), err.Error())
		}

		// no stream support currently
		// TODO: lookup host
		hreq, err := http.NewRequest("POST", ep, bytes.NewReader(body))
		if err != nil {
			return errors.InternalServerError(req.Service(), err.Error())
		}

		// get the header
		hdr := req.Header()

		// set the headers
		for k, v := range hdr {
			hreq.Header.Set(k, v)
		}

		// make the call
		hrsp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			return errors.InternalServerError(req.Service(), err.Error())
		}

		// read body
		b, err := ioutil.ReadAll(hrsp.Body)
		hrsp.Body.Close()
		if err != nil {
			return errors.InternalServerError(req.Service(), err.Error())
		}

		// set response headers
		hdr = map[string]string{}
		for k, _ := range hrsp.Header {
			hdr[k] = hrsp.Header.Get(k)
		}
		// write the header
		rsp.WriteHeader(hdr)
		// write the body
		err = rsp.Write(b)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.InternalServerError(req.Service(), err.Error())
		}
	}

	return nil
}

// NewSingleHostRouter returns a router which sends requests a single http backend
//
// It is used by setting it in a new go-micro to act as a proxy for a http backend.
//
// Usage:
//
// Create a new router to the http backend
//
// 	r := NewSingleHostRouter("http://localhost:10001")
//
// 	// Create your new service
// 	service := micro.NewService(
// 		micro.Name("greeter"),
// 	)
//
// 	// Setup the router
// 	service.Server().Init(server.WithRouter(r))
//
// 	// Run the service
// 	service.Run()
func NewSingleHostRouter(url string) *Router {
	return &Router{
		Resolver: new(Resolver),
		Endpoint: url,
		eps:      map[string]string{},
	}
}

// NewService returns a new http proxy. It acts as a go-micro service and proxies to http backend.
// Optionally specify the backend endpoint or the router. Otherwise a default is set.
//
// Usage:
//
// 	service := NewProxy(
//		micro.Name("greeter"),
//		http.SetRouter(r),
// 		// OR
//		http.SetEndpoint("http://localhost:10001"),
//	 )
func NewService(opts ...micro.Option) micro.Service {
	// prepend router to opts
	opts = append([]micro.Option{SetRouter(NewSingleHostRouter(Endpoint))}, opts...)

	// create the new service
	return emicro.NewService(opts...)
}
