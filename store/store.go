package store

import "github.com/noona-hq/app-blacklist/store/entity"

type Store interface {
	CreateUser(user entity.User) error
	UpdateUser(id string, user entity.User) (entity.User, error)
	GetUserForCompany(companyID string) (entity.User, error)
	DeleteUser(id string) error
}
