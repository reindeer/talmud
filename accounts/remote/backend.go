package remote

import (
	"github.com/reindeer/talmud/accounts/sqlite"
)

type Storage struct {
	sqlite.Storage
}

func New() *Storage {
	return &Storage{Storage: *sqlite.New()}
}
