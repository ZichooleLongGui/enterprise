// Package token is for tokenisation
package token

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hako/branca"
	p "github.com/micro/enterprise/proto"
	"github.com/pborman/uuid"
)

type Token struct {
	*p.Token
}

func (t *Token) Encode(key string) (string, error) {
	b, err := json.Marshal(t)
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
	t.Hash = str
	return str, nil
}

func (t *Token) Decode(key string, b []byte) error {
	if len(key) > 32 {
		key = key[:32]
	}
	br := branca.NewBranca(key)
	str, err := br.DecodeToString(string(b))
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(str), t); err != nil {
		return err
	}
	return nil
}

func (t *Token) Valid() error {
	// check id
	if len(t.Id) == 0 {
		return fmt.Errorf("token id invalid")
	}

	// check token expiry
	if u := time.Now().Unix(); (t.Expires - uint64(u)) < 0 {
		return fmt.Errorf("token expired")
	}

	// no claims
	if t.Claims == nil || len(t.Claims["email"]) == 0 {
		return fmt.Errorf("token claims invalid")
	}

	return nil
}

func New() *Token {
	return &Token{&p.Token{
		Id:      uuid.NewUUID().String(),
		Expires: uint64(time.Now().Add(time.Hour * 24 * 7).Unix()),
		Claims:  make(map[string]string),
	}}
}
