package main

import (
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/pborman/getopt"
	"github.com/reindeer/talmud/accounts"
	"github.com/reindeer/talmud/accounts/models"
	"github.com/reindeer/talmud/generator"
	"github.com/reindeer/talmud/master"
	"github.com/reindeer/talmud/output"
)

func main() {
	verbose := getopt.BoolLong("verbose", 'v', "Show account details")
	find := getopt.StringLong("find", 'f', "", "Find account by domain part")
	del := getopt.IntLong("delete", 'd', 0, "Delete account by Id")
	masterPassword := getopt.BoolLong("master-password", 'm', "Save master password to keychain")
	getopt.Parse()

	args := getopt.Args()
	if *del != 0 {
		accounts.Repository().Delete(*del)
	} else if masterPassword != nil && *masterPassword {
		master.New().Save()
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
			output.Fatalf("no accounts found")
		} else if len(list) == 1 {
			out(list[0], *verbose || find != nil)
		} else {
			showAccounts(list, *verbose)
		}
	} else {
		output.Fatalf("unknown request")
	}
}

func out(account *models.Account, verbose bool) {
	if verbose {
		showSingleAccount(account, verbose)
	}

	masterPassword := master.New().Get()

	gen := generator.New(&masterPassword)
	password, err := gen.Generate(account)
	if err != nil {
		output.Fatalf("cannot generate the password")
	}
	err = clipboard.WriteAll(*password)
	if err != nil {
		output.Printf(*password)
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

	output.Printf(mask, params...)
}
