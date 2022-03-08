package sqlite

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kirsle/configdir"
	_ "github.com/mattn/go-sqlite3"
	"github.com/reindeer/talmud/output"

	accounts "github.com/reindeer/talmud/accounts/models"
)

var schema = `
create table if not exists accounts (
    id integer primary key,
	domain varchar(255) not null,
	account varchar(255) not null,
	version int not null default 1,
	length int null,
	unique (domain, account)
)
`

type Storage struct {
	Db *sqlx.DB
}

func New() *Storage {
	path := getPath()
	db, err := sqlx.Open("sqlite3", path+"/pass.db")
	if err != nil {
		panic(err)
	}

	db.MustExec(schema)

	return &Storage{Db: db}
}

func (s Storage) List(domain *string) []*accounts.Account {
	query := `select * from accounts where version>0 order by domain, account`

	var list []*accounts.Account
	err := s.Db.Select(&list, query)
	if err != nil {
		panic(err)
	}

	var filtered []*accounts.Account
	for idx, account := range list {
		account.Idx = idx + 1
		if domain != nil && !strings.Contains(account.Domain, *domain) {
			continue
		}
		filtered = append(filtered, account)
	}

	return filtered
}

func (s Storage) Get(accountId int) *accounts.Account {
	list := s.List(nil)
	if accountId > len(list) {
		output.Fatalf("unknown account number #%d", accountId)
	}
	return list[accountId-1]
}

func (s Storage) Save(account *accounts.Account) {
	_, _ = s.Db.Exec(
		`insert into accounts (domain, account, version, length) values (?, ?, ?, ?) on conflict (domain, account) do update set version=?, length=?`,
		account.Domain, account.Account, account.Version, account.Length, account.Version, account.Length,
	)
}

func (s Storage) Delete(accountId int) {
	list := s.List(nil)
	account := list[accountId-1]
	_, _ = s.Db.Exec("update accounts set version=-1 where id=?", account.Id)
}

func getPath() string {
	path := configdir.LocalConfig("talmud")

	err := configdir.MakePath(path)
	if err != nil {
		panic(err)
	}

	return path
}
