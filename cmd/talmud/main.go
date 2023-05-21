package main

import (
	"fmt"
	"strconv"

	"github.com/reindeer/talmud/internal/adapter/interactor"
	"github.com/reindeer/talmud/internal/domain/account/model"
	accountManager "github.com/reindeer/talmud/internal/domain/account/service"
	masterManager "github.com/reindeer/talmud/internal/domain/master/service"
	"github.com/reindeer/talmud/pkg/try"

	"github.com/pborman/getopt"
)

func main() {
	var terminal interactor.Interactor = interactor.NewInteractor()
	accountManagement := accountManager.NewAccountManagement()
	masterPasswordManagement := masterManager.NewMasterPasswordManagement()

	try.Catch(
		func() {
			verbose := getopt.BoolLong("verbose", 'v', "Show account details")
			find := getopt.StringLong("find", 'f', "", "Find account by domain part")
			del := getopt.IntLong("delete", 'd', 0, "Delete account by Id")
			masterPassword := getopt.BoolLong("master-password", 'm', "Save master password to keychain")

			getopt.Parse()
			args := getopt.Args()

			if *del != 0 {
				accountManagement.Delete(*del)
			} else if masterPassword != nil && *masterPassword {
				password, _ := terminal.LoadPassword()
				masterPasswordManagement.Save(password)
			} else if len(args) == 1 {
				accountId, _ := strconv.Atoi(args[0])
				if (accountId) != 0 {
					account := accountManagement.Get(accountId)
					password := accountManagement.Password(account, masterPasswordManagement)
					terminal.PrintAccounts(*verbose, account)
					terminal.PrintPassword(password)
				}
			} else if len(args) >= 2 && len(args) <= 4 {
				args = append(args, "0", "0")

				version, _ := strconv.Atoi(args[2])
				if version == 0 {
					version = 1
				}

				length, _ := strconv.Atoi(args[3])

				account := accountManagement.Save(&model.Account{
					Domain:  args[0],
					Account: args[1],
					Version: version,
					Length:  length,
				})
				password := accountManagement.Password(account, masterPasswordManagement)
				terminal.PrintAccounts(*verbose, account)
				terminal.PrintPassword(password)
			} else if len(args) == 0 {
				accounts := accountManagement.List(find)
				if len(accounts) == 0 {
					panic(fmt.Errorf("no accounts found"))
				} else if len(accounts) == 1 {
					password := accountManagement.Password(accounts[0], masterPasswordManagement)
					terminal.PrintAccounts(*verbose || find != nil, accounts[0])
					terminal.PrintPassword(password)
				} else {
					terminal.PrintAccounts(*verbose, accounts...)
				}
			} else {
				panic(fmt.Errorf("unknown request"))
			}
		},
		func(throwable error) {
			terminal.Alert(throwable)
		},
	)
}
