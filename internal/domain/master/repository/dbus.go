//go:build dbus && !darwin && !ios

package repository

import (
	"fmt"

	"github.com/keybase/go-keychain/secretservice"
	dbus "github.com/keybase/go.dbus"

	"github.com/reindeer/talmud/pkg/try"
)

func NewMasterRepository() *dBus {
	return &dBus{
		label:      "talmud",
		attributes: map[string]string{"talmud": "master password"},
	}
}

type dBus struct {
	label      string
	attributes map[string]string
}

func (d *dBus) LoadPassword() string {
	srv := try.Throw(secretservice.NewService())
	session := try.Throw(srv.OpenSession(secretservice.AuthenticationDHAES))
	defer srv.CloseSession(session)

	collection := secretservice.DefaultCollection
	results := try.Throw(srv.SearchCollection(collection, d.attributes))
	if len(results) != 1 {
		panic(fmt.Errorf("%w: expected 1, %d given", ErrItemNotFound, len(results)))
	}

	secretPlaintext := try.Throw(srv.GetSecret(results[0], *session))

	return string(secretPlaintext)
}

func (d *dBus) SavePassword(password string) {
	srv := try.Throw(secretservice.NewService())
	session := try.Throw(srv.OpenSession(secretservice.AuthenticationDHAES))
	defer srv.CloseSession(session)

	collection := secretservice.DefaultCollection
	secret := try.Throw(session.NewSecret([]byte(password)))
	try.ThrowError(srv.Unlock([]dbus.ObjectPath{collection}))

	_ = try.Throw(srv.CreateItem(collection, secretservice.NewSecretProperties(d.label, d.attributes), secret, secretservice.ReplaceBehaviorReplace))
}
