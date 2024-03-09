package store

import "github.com/noona-hq/blacklist/store/entity"

type Store interface {
	CreateUser(user entity.User) error
	UpdateUser(id string, user entity.User) (entity.User, error)
	GetUserForCompany(companyID string) (entity.User, error)
}
