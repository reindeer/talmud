package interactor

import (
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/atotto/clipboard"
	"golang.org/x/term"

	"github.com/reindeer/talmud/internal/domain/account/model"
	"github.com/reindeer/talmud/pkg/try"
)

var (
	interactorInstance *interactor
	interactorFactory  sync.Once
)

func NewInteractor() *interactor {
	interactorFactory.Do(func() {
		interactorInstance = &interactor{}
	})
	return interactorInstance
}

type interactor struct{}

func (i *interactor) PrintAccounts(verbose bool, accounts ...*model.Account) {
	for _, account := range accounts {
		var params []any
		mask := "\033[1m%4d\033[0m: \033[33m%s\033[0m %s"
		params = append(params, account.Idx, account.Domain, account.Account)

		if account.Length != 0 {
			mask += " (%d symbols)"
			params = append(params, account.Length)
		} else if verbose {
			mask += " (no limitations)"
		}
		if verbose {
			mask += " version %d"
			params = append(params, account.Version)
		}

		fmt.Printf(mask+"\n", params...)
	}
}

func (i *interactor) LoadPassword() (string, bool) {
	fmt.Printf("Enter master password: ")

	password := try.Throw(term.ReadPassword(syscall.Stdin))
	fmt.Println("")

	return string(password), false
}

func (i *interactor) PrintPassword(password string) {
	if err := clipboard.WriteAll(password); err != nil {
		fmt.Printf("%s\n", password)
	}
}

func (i *interactor) Alert(err error) {
	fmt.Printf("\033[31m%v\033[0m\n", err)
	os.Exit(1)
}
