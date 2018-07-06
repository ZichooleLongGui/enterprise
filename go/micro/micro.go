package micro

import (
	"github.com/micro/go-micro"
)

// NewService returns a new enterprise go-micro Service
func NewService(opts ...micro.Option) micro.Service {
	return micro.NewService(opts...)
}
