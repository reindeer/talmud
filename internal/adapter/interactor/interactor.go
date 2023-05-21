package interactor

import (
	"github.com/reindeer/talmud/internal/domain/account/model"
)

type Interactor interface {
	PrintAccounts(verbose bool, accounts ...*model.Account)
	LoadPassword() (string, bool)
	PrintPassword(password string)
	Alert(err error)
}
