//go:build !darwin && !ios && !dbus

package repository

import "sync"

var (
	repositoryInstance *generic
	repositoryFactory  sync.Once
)

func NewMasterRepository() *generic {
	repositoryFactory.Do(func() {
		repositoryInstance = &generic{}
	})
	return repositoryInstance
}

type generic struct {
}

func (i *generic) LoadPassword() string {
	return ""
}

func (i *generic) SavePassword(string) {
}
