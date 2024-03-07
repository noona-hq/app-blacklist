package mongodb

import (
	"context"
	"time"

	"github.com/noona-hq/blacklist/db"
	"github.com/noona-hq/blacklist/services/store"
	"github.com/noona-hq/blacklist/services/store/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	blacklistCollectionName = "blacklist"
)

type Store struct {
	db db.Database
}

// MongoDB implementation for store
func NewStore(db db.Database) store.Store {
	return Store{
		db: db,
	}
}

func (s Store) CreateBlacklistUser(user entity.User) error {
	blacklistCollection := s.db.DB.Collection(blacklistCollectionName)

	if user.ID == "" {
		user.ID = randomID()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := blacklistCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) GetBlacklistUserForCompany(companyID string) (entity.User, error) {
	blacklistCollection := s.db.DB.Collection(blacklistCollectionName)

	filter := filter()

	filter["companyID"] = companyID

	// Sort by createdAt descending
	sort := bson.M{"createdAt": -1}

	var user entity.User
	err := blacklistCollection.FindOne(context.Background(), filter, &options.FindOneOptions{Sort: sort}).Decode(&user)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func filter() bson.M {
	filter := bson.M{"deletedAt": bson.M{"$exists": false}}

	return filter
}
