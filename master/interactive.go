package master

import (
	"fmt"
	"syscall"

	"github.com/reindeer/talmud/output"
	"golang.org/x/term"
)

type Interactive interface {
	Get() string
	Save() string
}

type Console struct {
}

func (Console) Get() string {
	return scan()
}

func (Console) Save() string {
	return scan()
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
