//go:build darwin || ios

package repository

import (
	"fmt"
	"sync"

	"github.com/keybase/go-keychain"

	"github.com/reindeer/talmud/pkg/try"
)

var (
	repositoryInstance *keyChain
	repositoryFactory  sync.Once
)

func NewMasterRepository() *keyChain {
	repositoryFactory.Do(func() {
		repositoryInstance = &keyChain{account: "talmud"}
	})
	return repositoryInstance
}

type keyChain struct {
	account string
}

func (k *keyChain) LoadPassword() string {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetAccount(k.account)
	query.SetReturnAttributes(true)
	query.SetReturnData(true)
	results := try.Throw(keychain.QueryItem(query))

	if len(results) != 1 {
		panic(fmt.Errorf("%w: expected 1, %d given", ErrItemNotFound, len(results)))
	}

	return string(results[0].Data)
}

func (k *keyChain) SavePassword(password string) {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetAccount(k.account)
	item.SetData([]byte(password))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	try.ThrowError(keychain.DeleteItem(item))
	try.ThrowError(keychain.AddItem(item))
}
