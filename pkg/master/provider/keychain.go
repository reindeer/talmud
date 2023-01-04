//go:build darwin || ios

package provider

import (
	"fmt"
	"github.com/keybase/go-keychain"
	"github.com/reindeer/talmud/internal/terminal"
	"github.com/reindeer/talmud/internal/try"
)

func Get() string {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetAccount("talmud")
	query.SetReturnAttributes(true)
	query.SetReturnData(true)
	results := try.Throw(keychain.QueryItem(query))

	if len(results) != 1 {
		panic(fmt.Errorf("expected 1 master password, %d given", len(results)))
	}

	return string(results[0].Data)
}

func Save() {
	password := terminal.Scan()
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetAccount("talmud")
	item.SetData([]byte(password))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	try.ThrowError(keychain.DeleteItem(item))
	try.ThrowError(keychain.AddItem(item))
}
