// +build darwin ios

package master

import (
	"github.com/keybase/go-keychain"
	"github.com/reindeer/talmud/output"
)

func get() string {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetAccount("talmud")
	query.SetReturnAttributes(true)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		output.Fatalf("%v", err)
	} else if len(results) != 1 {
		output.Fatalf("expected 1 master password, %d given", len(results))
	}

	return string(results[0].Data)
}

func save() {
	password := scan()
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetAccount("talmud")
	item.SetData([]byte(password))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	err := keychain.DeleteItem(item)
	if err != nil {
		output.Fatalf("cannot delete the old master password: %v", err)
	}

	err = keychain.AddItem(item)
	if err == keychain.ErrorDuplicateItem {
		output.Fatalf("cannot save the master password")
	}
}
