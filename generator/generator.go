package generator

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"

	accounts "github.com/reindeer/talmud/accounts/models"
)

type Generator struct {
	MP *string
}

func New(masterPassword *string) *Generator {
	return &Generator{MP: masterPassword}
}

func (g Generator) Generate(account *accounts.Account) (*string, error) {
	if g.MP == nil {
		return nil, errors.New("empty master password")
	}

	masterPassword := encrypt(sha256.New(), []byte(*g.MP))

	login := fmt.Sprintf("%s:%s:%d", account.Domain, account.Account, account.Version)
	hmacPassword := encrypt(hmac.New(sha256.New, masterPassword), []byte(login))

	password := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(encrypt(sha1.New(), hmacPassword))
	if account.Length > 0 && account.Length < len(password) {
		password = password[:account.Length]
	}
	return &password, nil
}

func encrypt(hasher hash.Hash, message []byte) []byte {
	hasher.Write(message)
	return hasher.Sum(nil)
}
