package http

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
)

// SetEndpoint provides an option to set the http backend url
func SetEndpoint(url string) micro.Option {
	return func(o *micro.Options) {
		// get the router
		r := o.Server.Options().Router
		// check its a http router
		if httpRouter, ok := r.(*Router); ok {
			httpRouter.Endpoint = url
		}
	}
}

// SetRouter provides an option to set the http router
func SetRouter(r server.Router) micro.Option {
	return func(o *micro.Options) {
		o.Server.Init(server.WithRouter(r))
	}
}
