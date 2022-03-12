//go:build dbus && !darwin && !ios
// +build dbus,!darwin,!ios

package master

import (
	"github.com/keybase/go-keychain/secretservice"
	dbus "github.com/keybase/go.dbus"
	"github.com/reindeer/talmud/output"
)

func get() string {
	srv, err := secretservice.NewService()
	if err != nil {
		output.Fatalf("cannot get the master password: %v", err)
	}

	session, err := srv.OpenSession(secretservice.AuthenticationDHAES)
	if err != nil {
		output.Fatalf("cannot get the master password: %v", err)
	}
	defer srv.CloseSession(session)

	collection := secretservice.DefaultCollection
	items, err := srv.SearchCollection(collection, map[string]string{"talmud": "master password"})
	if err != nil {
		output.Fatalf("%v", err)
	} else if len(items) != 1 {
		output.Fatalf("expected 1 master password, %d given", len(items))
	}

	secretPlaintext, err := srv.GetSecret(items[0], *session)
	if err != nil {
		output.Fatalf("cannot get the master password: %v", err)
	}

	return string(secretPlaintext)
}

func save() {
	password := scan()

	srv, err := secretservice.NewService()
	if err != nil {
		output.Fatalf("cannot save the master password: %v", err)
	}

	session, err := srv.OpenSession(secretservice.AuthenticationDHAES)
	if err != nil {
		output.Fatalf("cannot save the master password: %v", err)
	}
	defer srv.CloseSession(session)

	collection := secretservice.DefaultCollection

	secret, err := session.NewSecret([]byte(password))
	if err != nil {
		output.Fatalf("cannot save the master password: %v", err)
	}

	err = srv.Unlock([]dbus.ObjectPath{collection})
	if err != nil {
		output.Fatalf("cannot save the master password: %v", err)
	}

	_, err = srv.CreateItem(collection, secretservice.NewSecretProperties("talmud", map[string]string{"talmud": "master password"}), secret, secretservice.ReplaceBehaviorReplace)
	if err != nil {
		output.Fatalf("cannot save the master password: %v", err)
	}
}
