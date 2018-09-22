// Package license provides license management code
package license

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	u = "https://micro.mu/license/"

	// license version
	v = "20180905"

	// api token
	t = os.Getenv("MICRO_TOKEN_KEY")

	// enterprise license
	l = os.Getenv("MICRO_LICENSE_KEY")
)

func init() {
	if uri := os.Getenv("MICRO_LICENSE_API"); len(uri) > 0 {
		u = uri
	}
}

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

// Info is the response from the update
type Info struct {
	*proto.Info
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
	if (l.Subscription.Type != "service") && (l.Subscription.Type != "unlimited") {
		return fmt.Errorf("subscription type invalid")
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
	if len(u.Id) == 0 {
		return fmt.Errorf("update id is invalid")
	}
	if u.Timestamp < uint64(c.Unix()) {
		return fmt.Errorf("update timestamp is invalid")
	}
	if u.Service == nil {
		return fmt.Errorf("update service is invalid")
	}

	us := &Service{u.Service}

	if err := us.Valid(); err != nil {
		return err
	}
	return nil
}

func call(method, uri string, vals url.Values) (*http.Response, error) {
	// check token
	if len(t) == 0 {
		return nil, fmt.Errorf("Require MICRO_TOKEN_KEY")
	}
	// set vals
	var data io.Reader
	if vals != nil {
		data = strings.NewReader(vals.Encode())
	}
	req, err := http.NewRequest(method, uri, data)
	if err != nil {
		return nil, err
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("X-Micro-Token", t)
	return http.DefaultClient.Do(req)
}

// SendUpdate sends a license update
func SendUpdate(ud *Update) (*Info, error) {
	b, err := json.Marshal(ud)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", u+"update", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Micro-License", l)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	b, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Api error: %s (require MICRO_LICENSE_KEY)", strings.TrimSpace(string(b)))
	}
	var info *Info
	if err := json.Unmarshal(b, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// Generate generates the license
func Generate(subscription string) (string, error) {
	data := url.Values{
		"subscription": {subscription},
	}
	rsp, err := call("POST", u+"generate", data)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != 200 {
		return "", fmt.Errorf(string(b))
	}
	var res map[string]interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		return "", err
	}
	license, _ := res["license"].(string)
	return license, nil
}

// Revoke revokes a license
func Revoke(lu string) error {
	data := url.Values{
		"license": {lu},
	}
	rsp, err := call("POST", u+"revoke", data)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode == 401 {
		return fmt.Errorf("Api error: %s (require MICRO_TOKEN_KEY)", strings.TrimSpace(string(b)))
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf("API error: %s", strings.TrimSpace(string(b)))
	}
	return nil
}

// List lists the licenses
func List() ([]*License, error) {
	rsp, err := call("GET", u+"list", nil)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode == 401 {
		return nil, fmt.Errorf("Api error: %s (require MICRO_TOKEN_KEY)", strings.TrimSpace(string(b)))
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", strings.TrimSpace(string(b)))
	}
	var list map[string][]*License
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	return list["licenses"], nil
}

// Verify a token is valid
func Verify(lu string) error {
	data := url.Values{
		"license": {lu},
	}
	rsp, err := call("POST", u+"verify", data)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode == 401 {
		return fmt.Errorf(strings.TrimSpace(string(b)))
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf(strings.TrimSpace(string(b)))
	}
	return nil
}

// Subscriptions lists the subscriptions
func Subscriptions() ([]*Subscription, error) {
	rsp, err := call("GET", u+"subscriptions", nil)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode == 401 {
		return nil, fmt.Errorf("Api error: %s (require MICRO_TOKEN_KEY)", strings.TrimSpace(string(b)))
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", strings.TrimSpace(string(b)))
	}
	var list []*Subscription
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// SetToken sets the api token
func SetToken(tk string) {
	t = tk
}

// Set license to use on update calls
func SetLicense(lu string) {
	l = lu
}

func NewUpdate() *Update {
	return &Update{&proto.Update{
		Id:        uuid.NewUUID().String(),
		Timestamp: uint64(time.Now().Unix()),
		Service:   &proto.Service{},
	}}
}

func New() *License {
	return &License{&proto.License{
		Id:           uuid.NewUUID().String(),
		Version:      v,
		Created:      uint64(time.Now().Unix()),
		Subscription: &proto.Subscription{},
	}}
}
