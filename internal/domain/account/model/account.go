package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"time"
)

type Account struct {
	Id        int        `json:"id,omitempty" db:"id,omitempty"`
	Idx       int        `json:"-" db:"-"`
	Domain    string     `json:"domain" db:"domain"`
	Account   string     `json:"account" db:"account"`
	Version   int        `json:"version" db:"version"`
	Length    int        `json:"length" db:"length"`
	Deleted   bool       `json:"deleted" db:"deleted"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	Comment   *string    `json:"comment,omitempty" db:"comment,omitempty"`
}

func (a *Account) Password(master string) string {
	masterPassword := a.encrypt(sha256.New(), []byte(master))
	hmacPassword := a.encrypt(hmac.New(sha256.New, masterPassword), []byte(fmt.Sprintf("%s:%s:%d", a.Domain, a.Account, a.Version)))

	password := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(a.encrypt(sha1.New(), hmacPassword))
	if a.Length > 0 && a.Length < len(password) {
		return password[:a.Length]
	}
	return password
}

func (a *Account) encrypt(hasher hash.Hash, message []byte) []byte {
	hasher.Write(message)
	return hasher.Sum(nil)
}
