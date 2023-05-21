package db

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/kirsle/configdir"

	"github.com/reindeer/talmud/pkg/try"
)

var schema = `
create table if not exists accounts (
    id integer primary key,
	domain varchar(255) not null,
	account varchar(255) not null,
	version int not null default 1,
	length int null,
	deleted int(1) default 0,
	created_at datetime null,
	updated_at datetime null,
	comment text null,
	unique (domain, account)
);

create trigger if not exists [UpdatedAt]
after insert on accounts
for each row 
begin
	update accounts set created_at=CURRENT_TIMESTAMP where id=new.id;
end;

create trigger if not exists [CreatedAt]
after update on accounts
for each row 
begin
	update accounts set updated_at=CURRENT_TIMESTAMP where id=old.id;
end;

`

var (
	client        *sqlx.DB
	clientFactory sync.Once
)

func NewConnection() *sqlx.DB {
	clientFactory.Do(func() {
		client = try.Throw(sqlx.Open("sqlite3", getPath()+"/pass.db"))
		client.MustExec(schema)
	})

	return client
}

func getPath() string {
	path := configdir.LocalConfig("talmud")
	try.ThrowError(configdir.MakePath(path))
	return path
}
