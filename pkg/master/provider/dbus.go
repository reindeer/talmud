//go:build dbus && !darwin && !ios

package provider

import (
	"fmt"
	"github.com/keybase/go-keychain/secretservice"
	dbus "github.com/keybase/go.dbus"
	"github.com/reindeer/talmud/internal/terminal"
	"github.com/reindeer/talmud/internal/try"
)

func Get() string {
	srv := try.Throw(secretservice.NewService())
	session := try.Throw(srv.OpenSession(secretservice.AuthenticationDHAES))
	defer srv.CloseSession(session)

	collection := secretservice.DefaultCollection
	items := try.Throw(srv.SearchCollection(collection, map[string]string{"talmud": "master password"}))
	if len(items) != 1 {
		panic(fmt.Errorf("expected 1 master password, %d given", len(items)))
	}

	secretPlaintext := try.Throw(srv.GetSecret(items[0], *session))

	return string(secretPlaintext)
}

func Save() {
	password := terminal.Scan()

	srv := try.Throw(secretservice.NewService())

	session := try.Throw(srv.OpenSession(secretservice.AuthenticationDHAES))
	defer srv.CloseSession(session)

	collection := secretservice.DefaultCollection

	secret := try.Throw(session.NewSecret([]byte(password)))
	try.ThrowError(srv.Unlock([]dbus.ObjectPath{collection}))

	_ = try.Throw(srv.CreateItem(collection, secretservice.NewSecretProperties("talmud", map[string]string{"talmud": "master password"}), secret, secretservice.ReplaceBehaviorReplace))
}
