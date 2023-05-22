package repository

import (
	"sync"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"

	"github.com/reindeer/talmud/internal/adapter/db"
	"github.com/reindeer/talmud/internal/domain/account/model"
	"github.com/reindeer/talmud/pkg/try"
)

type Account interface {
	All() []*model.Account
	Save(account *model.Account) *model.Account
	Delete(account *model.Account)
}

var (
	repositoryInstance *account
	repositoryFactory  sync.Once
)

func NewAccountRepository() *account {
	repositoryFactory.Do(func() {
		repositoryInstance = &account{db: db.NewConnection()}
	})
	return repositoryInstance
}

type account struct {
	db *sqlx.DB
}

func (a *account) All() []*model.Account {
	query := `select * from accounts where deleted=0 order by domain, account`

	var accounts []*model.Account
	try.ThrowError(a.db.Select(&accounts, query))

	for idx := range accounts {
		accounts[idx].Idx = idx + 1
	}
	return accounts
}

func (a *account) Save(account *model.Account) *model.Account {
	result := try.Throw(a.db.Exec(
		`insert into accounts (domain, account, version, length) values (?, ?, ?, ?) on conflict (domain, account) do update set version=?, length=?, deleted=0`,
		account.Domain, account.Account, account.Version, account.Length, account.Version, account.Length,
	))
	account.Id = int(try.Throw(result.LastInsertId()))
	return account
}

func (a *account) Delete(account *model.Account) {
	_, _ = a.db.Exec("update accounts set deleted=1 where id=?", account.Id)
}
