package repository

import (
	"fmt"
)

var (
	ErrItemNotFound = fmt.Errorf("cannot load item")
)

type Repository interface {
	LoadPassword() string
	SavePassword(password string)
}
