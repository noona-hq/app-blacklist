package store

import "github.com/noona-hq/blacklist/services/store/entity"

type Store interface {
	CreateBlacklistUser(user entity.User) error
	UpdateBlacklistUser(id string, user entity.User) (entity.User, error)
	GetBlacklistUserForCompany(companyID string) (entity.User, error)
}
