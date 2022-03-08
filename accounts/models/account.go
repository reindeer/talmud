package models

type Account struct {
	Id      int    `json:"id,omitempty"`
	Idx     int    `json:"-"`
	Domain  string `json:"domain"`
	Account string `json:"account"`
	Version int    `json:"version"`
	Length  int    `json:"length"`
}
