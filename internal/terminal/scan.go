package terminal

import (
	"fmt"
	"github.com/reindeer/talmud/internal/try"
	"golang.org/x/term"
	"syscall"
)

func Scan() string {
	fmt.Printf("Enter master password: ")

	password := try.Throw(term.ReadPassword(syscall.Stdin))
	fmt.Println("")

	return string(password)
}
