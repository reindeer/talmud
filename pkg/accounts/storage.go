package accounts

import (
	"github.com/reindeer/talmud/pkg/accounts/models"
	"github.com/reindeer/talmud/pkg/accounts/sqlite"
)

type RepositoryInterface interface {
	List(domain *string) []*models.Account
	Get(accountId int) *models.Account
	Save(account *models.Account)
	Delete(accountId int)
}

type Storage struct {
	RepositoryInterface
}

var repository *Storage

func Repository() *Storage {
	if repository == nil {
		storage := sqlite.New()
		repository = &Storage{RepositoryInterface: storage}
	}
	return repository
}
