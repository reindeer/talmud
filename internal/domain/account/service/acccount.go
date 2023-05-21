package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/reindeer/talmud/internal/domain/account/model"
	"github.com/reindeer/talmud/internal/domain/account/repository"
	"github.com/reindeer/talmud/internal/domain/master/service"
	"github.com/reindeer/talmud/pkg/list"
)

var (
	ErrAccountNotFound     = fmt.Errorf("account is not found")
	ErrEmptyMasterPassword = fmt.Errorf("empty master password")
)

type AccountManagement interface {
	List(domain *string) []*model.Account
	Get(accountId int) *model.Account
	Password(account *model.Account, masterManagement service.MasterPasswordManagement) string
	Save(account *model.Account) *model.Account
	Delete(accountId int)
}

var (
	managementInstance *account
	managementFactory  sync.Once
)

func NewAccountManagement() *account {
	managementFactory.Do(func() {
		managementInstance = &account{repository: repository.NewAccountRepository()}
	})
	return managementInstance
}

type account struct {
	repository repository.Account
}

func (a *account) List(domain *string) []*model.Account {
	return list.Filter(
		a.repository.All(),
		func(account *model.Account) bool {
			return domain == nil || strings.Contains(account.Domain, *domain)
		},
	)
}

func (a *account) Get(accountId int) *model.Account {
	accounts := a.List(nil)
	if accountId > len(accounts) {
		panic(fmt.Errorf("%w: #%d", ErrAccountNotFound, accountId))
	}
	return accounts[accountId-1]
}

func (a *account) Password(account *model.Account, masterManagement service.MasterPasswordManagement) string {
	masterPassword := masterManagement.Load()
	if masterPassword == "" {
		panic(ErrEmptyMasterPassword)
	}
	return account.Password(masterPassword)
}

func (a *account) Save(account *model.Account) *model.Account {
	return a.repository.Save(account)
}

func (a *account) Delete(accountId int) {
	account := a.Get(accountId)
	a.repository.Delete(account)
}
