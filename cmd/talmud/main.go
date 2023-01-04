package main

import (
	"fmt"
	"github.com/reindeer/talmud/internal/try"
	"os"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/pborman/getopt"
	"github.com/reindeer/talmud/pkg/accounts"
	"github.com/reindeer/talmud/pkg/accounts/models"
	"github.com/reindeer/talmud/pkg/generator"
	"github.com/reindeer/talmud/pkg/master"
)

func main() {
	try.Catch(
		func() {
			exec()
		},
		func(throwable error) {
			fmt.Printf("\033[31m%v\033[0m\n", throwable)
			os.Exit(1)
		},
	)
}

func exec() {
	verbose := getopt.BoolLong("verbose", 'v', "Show account details")
	find := getopt.StringLong("find", 'f', "", "Find account by domain part")
	del := getopt.IntLong("delete", 'd', 0, "Delete account by Id")
	masterPassword := getopt.BoolLong("master-password", 'm', "Save master password to keychain")

	getopt.Parse()
	args := getopt.Args()

	if *del != 0 {
		accounts.Repository().Delete(*del)
	} else if masterPassword != nil && *masterPassword {
		master.Save()
	} else if len(args) == 1 {
		accountId, _ := strconv.Atoi(args[0])
		if (accountId) != 0 {
			account := accounts.Repository().Get(accountId)
			out(account, *verbose)
		}
	} else if len(args) >= 2 && len(args) <= 4 {
		args = append(args, "0", "0")

		version, _ := strconv.Atoi(args[2])
		if version == 0 {
			version = 1
		}

		length, _ := strconv.Atoi(args[3])

		account := &models.Account{
			Domain:  args[0],
			Account: args[1],
			Version: version,
			Length:  length,
		}
		accounts.Repository().Save(account)
		out(account, *verbose)
	} else if len(args) == 0 {
		list := accounts.Repository().List(find)
		if len(list) == 0 {
			panic(fmt.Errorf("no accounts found"))
		} else if len(list) == 1 {
			out(list[0], *verbose || find != nil)
		} else {
			showAccounts(list, *verbose)
		}
	} else {
		panic(fmt.Errorf("unknown request"))
	}
}
func out(account *models.Account, verbose bool) {
	if verbose {
		showSingleAccount(account, verbose)
	}

	masterPassword := master.Get()

	gen := generator.NewGenerator(&masterPassword)
	password := gen.Generate(account.Domain, account.Account, account.Version, account.Length)
	if err := clipboard.WriteAll(password); err != nil {
		fmt.Printf("%s\n", password)
	}
}

func showAccounts(accounts []*models.Account, verbose bool) {
	for _, account := range accounts {
		showSingleAccount(account, verbose)
	}
}

func showSingleAccount(account *models.Account, verbose bool) {
	var params []interface{}
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
