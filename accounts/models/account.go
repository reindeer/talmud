package models

import "time"

type Account struct {
	Id        int        `json:"id,omitempty" db:"id,omitempty"`
	Idx       int        `json:"-" db:"-"`
	Domain    string     `json:"domain" db:"domain"`
	Account   string     `json:"account" db:"account"`
	Version   int        `json:"version" db:"version"`
	Length    int        `json:"length" db:"length"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	Comment   *string    `json:"comment,omitempty" db:"comment,omitempty"`
}
