package master

import (
	"fmt"
	"syscall"

	"github.com/reindeer/talmud/output"
	"golang.org/x/term"
)

func Get() string {
	return get()
}

func Save() {
	save()
}

func scan() string {
	fmt.Printf("Enter master password: ")

	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		output.Fatalf("%v", err)
	}
	output.Printf("")
	return string(password)
}
