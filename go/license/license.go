// package license provides related code
package license

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hako/branca"
	proto "github.com/micro/enterprise/proto"
	"github.com/pborman/uuid"
)

var (
	// first valid license date
	// 5 september 2018
	c = time.Date(2018, time.September, 5, 16, 23, 0, 0, time.UTC)

	// license api
	u = "http://localhost:9091/"

	// license version
	v = "20180905"
)

// License is the enterprise license
type License struct {
	*proto.License
}

// Service uses the license
type Service struct {
	*proto.Service
}

// Subscription is a user subscription
type Subscription struct {
	*proto.Subscription
}

// Update represents an update
type Update struct {
	*proto.Update
}

func (l *License) Encode(key string) (string, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return "", err
	}
	if len(key) > 32 {
		key = key[:32]
	}
	br := branca.NewBranca(key)
	str, err := br.EncodeToString(string(b))
	if err != nil {
		return "", err
	}
	l.Key = str
	return str, nil
}

func (l *License) Decode(key string, b []byte) error {
	if len(key) > 32 {
		key = key[:32]
	}
	br := branca.NewBranca(key)
	str, err := br.DecodeToString(string(b))
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(str), l); err != nil {
		return err
	}
	return nil
}

func (l *License) Equal(lu *License) error {
	if l.Id != lu.Id {
		return fmt.Errorf("invalid license id")
	}
	if l.Version != lu.Version {
		return fmt.Errorf("invalid license version")
	}
	su1 := &Subscription{lu.Subscription}
	su2 := &Subscription{l.Subscription}
	return su1.Equal(su2)
}

func (l *License) Valid() error {
	str := "license %s is blank"

	if len(l.Id) == 0 {
		return fmt.Errorf(str, "id")
	}
	if len(l.Version) == 0 {
		return fmt.Errorf(str, "version")
	}
	if l.Created < uint64(c.Unix()) {
		return fmt.Errorf("license creation time %d is invalid", l.Created)
	}
	if l.Subscription == nil {
		return fmt.Errorf("license subscription is nil")
	}
	if len(l.Subscription.Id) == 0 {
		return fmt.Errorf(str, "subscription id")
	}
	if len(l.Subscription.Email) == 0 {
		return fmt.Errorf(str, "subscription email")
	}
	if l.Subscription.Created < uint64(c.Unix()) {
		return fmt.Errorf("license subscription time %d is invalid", l.Subscription.Created)
	}
	return nil
}

// micro://email/subscription_id/license_id
func (l *License) String() string {
	return fmt.Sprintf("micro://%s/%s/%s", l.Subscription.Email, l.Subscription.Id, l.Id)
}

func (s *Subscription) Equal(su *Subscription) error {
	// email match
	if su.Email != s.Email {
		return fmt.Errorf("Email does not match subscription")
	}
	// id match
	if su.Id != s.Id {
		return fmt.Errorf("Id does not match subscription")
	}
	return nil
}

// micro://email/subscription
func (s *Subscription) String() string {
	return fmt.Sprintf("micro://%s/%s", s.Email, s.Id)
}

func (s *Service) Valid() error {
	str := "service %s is blank"

	if len(s.Name) == 0 {
		return fmt.Errorf(str, "name")
	}
	if len(s.Id) == 0 {
		return fmt.Errorf(str, "id")
	}
	if len(s.Version) == 0 {
		return fmt.Errorf(str, "version")
	}
	return nil
}

func (u *Update) Encode(key string) (string, error) {
	b, err := json.Marshal(u.Update)
	if err != nil {
		return "", err
	}
	if len(key) > 32 {
		key = key[:32]
	}
	br := branca.NewBranca(key)
	str, err := br.EncodeToString(string(b))
	if err != nil {
		return "", err
	}
	return str, nil
}

func (u *Update) Decode(key string, b []byte) error {
	if len(key) > 32 {
		key = key[:32]
	}
	br := branca.NewBranca(key)
	str, err := br.DecodeToString(string(b))
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(str), u.Update); err != nil {
		return err
	}
	return nil
}

func (u *Update) Valid() error {
	if u.Service == nil {
		return fmt.Errorf("service is nil")
	}
	if u.License == nil {
		return fmt.Errorf("license is nil")
	}
	if u.License.Subscription == nil {
		return fmt.Errorf("subscription is nil")
	}
	if u.Timestamp < uint64(c.Unix()) || u.Timestamp < u.License.Created {
		return fmt.Errorf("update timestamp is invalid")
	}
	if u.Timestamp < uint64(c.Unix()) || u.Timestamp < u.License.Created {
		return fmt.Errorf("update timestamp is invalid")
	}
	if len(u.Service.Name) == 0 {
		return fmt.Errorf("service name is blank")
	}
	if len(u.Service.Version) == 0 {
		return fmt.Errorf("service version is blank")
	}

	ul := &License{u.License}
	us := &Service{u.Service}

	if err := ul.Valid(); err != nil {
		return err
	}
	if err := us.Valid(); err != nil {
		return err
	}
	return nil
}

func New() *License {
	return &License{&proto.License{
		Id:           uuid.NewUUID().String(),
		Version:      v,
		Created:      uint64(time.Now().Unix()),
		Subscription: &proto.Subscription{},
	}}
}
