// Package micro is for enterprise Go Micro
package micro

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/micro/enterprise/go/license"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

type service struct {
	micro.Service
}

type licenseKey struct{}

// the update loop for sending updates
func (s *service) update(exit chan bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	next := time.Minute
	srv := s.Service.Server()
	name := srv.Options().Name
	version := srv.Options().Version
	id := srv.Options().Id

	backoff := func(attempt int) time.Duration {
		min := float64(time.Minute)
		max := float64(time.Hour * 24)

		d := min * math.Pow(10, float64(attempt))
		if d > max {
			d = max
		}
		return time.Duration(r.Float64()*(d-min) + min)
	}

	update := func(u *license.Update) *license.Update {
		if u != nil {
			return u
		}
		u = license.NewUpdate()
		u.Service.Name = name
		u.Service.Version = version
		u.Service.Id = id
		return u
	}

	// retries
	var i int

	// update
	var u *license.Update

	for {
		select {
		// exit
		case <-exit:
			return
		// update
		case <-time.After(next):
			// first attempt
			if i == 0 || u == nil {
				u = update(nil)
			}

			// send update
			info, err := license.SendUpdate(u)
			if err != nil {
				// backoff
				log.Logf("Sending update failed: %v", err)
				// set backoff
				next = backoff(i)
				// update counter
				i++
				continue
			}

			// next update
			next = time.Duration(info.NextUpdate-uint64(time.Now().Unix())) * time.Second

			// reset backoff
			i = 0
		}
	}
}

func (s *service) Run() error {
	// get license
	ctx := s.Service.Options().Context
	key, ok := ctx.Value(licenseKey{}).(string)
	if !ok {
		// try env var
		key = os.Getenv("MICRO_LICENSE_KEY")
	}

	// TODO: check key is valid
	if len(key) < 62 {
		return errors.New("micro enterprise license key missing")
	}

	// set the license
	license.SetLicense(key)

	exit := make(chan bool)
	go s.update(exit)

	defer func() {
		close(exit)
	}()

	return s.Service.Run()
}

// License is an option to set the license
func License(key string) micro.Option {
	return func(o *micro.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, licenseKey{}, key)
	}
}

// NewService returns a new enterprise Go Micro Service
func NewService(opts ...micro.Option) micro.Service {
	svc := &service{
		micro.NewService(opts...),
	}
	return svc
}
