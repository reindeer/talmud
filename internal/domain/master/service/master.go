package service

import (
	"sync"

	"github.com/reindeer/talmud/internal/adapter/interactor"
	"github.com/reindeer/talmud/internal/domain/master/repository"
)

type MasterPasswordManagement interface {
	Load() string
	Save(password string)
}

var (
	managementInstance *masterPassword
	managementFactory  sync.Once
)

func NewMasterPasswordManagement() *masterPassword {
	managementFactory.Do(func() {
		managementInstance = &masterPassword{
			provider:   repository.NewMasterRepository(),
			interactor: interactor.NewInteractor(),
		}
	})
	return managementInstance
}

type masterPassword struct {
	provider   repository.Repository
	interactor interactor.Interactor
}

func (m *masterPassword) Load() string {
	password := m.provider.LoadPassword()
	if password == "" {
		pass, save := m.interactor.LoadPassword()
		if save && pass != "" {
			m.provider.SavePassword(pass)
		}
		password = pass
	}
	return password
}

func (m *masterPassword) Save(password string) {
	m.provider.SavePassword(password)
}
