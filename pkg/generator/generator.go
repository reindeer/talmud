package generator

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
)

type Generator struct {
	MP *string
}

func NewGenerator(masterPassword *string) *Generator {
	return &Generator{MP: masterPassword}
}

func (g Generator) Generate(domain, account string, version, length int) string {
	login := fmt.Sprintf("%s:%s:%d", domain, account, version)
	password := g.GenerateFromString(login)

	if length > 0 && length < len(password) {
		return password[:length]
	}
	return password
}

func (g Generator) GenerateFromString(account string) string {
	if g.MP == nil {
		panic(errors.New("empty master password"))
	}

	masterPassword := encrypt(sha256.New(), []byte(*g.MP))
	hmacPassword := encrypt(hmac.New(sha256.New, masterPassword), []byte(account))

	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(encrypt(sha1.New(), hmacPassword))
}

func encrypt(hasher hash.Hash, message []byte) []byte {
	hasher.Write(message)
	return hasher.Sum(nil)
}
