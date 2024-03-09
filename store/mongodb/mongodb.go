package mongodb

import (
	"context"
	"time"

	"github.com/chidiwilliams/flatbson"
	"github.com/dchest/uniuri"
	"github.com/noona-hq/blacklist/db"
	"github.com/noona-hq/blacklist/store"
	"github.com/noona-hq/blacklist/store/entity"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersCollectionName = "users"
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

func (s Store) CreateUser(user entity.User) error {
	usersCollection := s.db.DB.Collection(usersCollectionName)

	if user.ID == "" {
		user.ID = randomID()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) UpdateUser(id string, user entity.User) (entity.User, error) {
	usersCollection := s.db.DB.Collection(usersCollectionName)

	user.UpdatedAt = time.Now()

	filter := filter()

	filter["_id"] = id

	userUpdate, err := flatbson.Flatten(user)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "Error flattening user")
	}

	update := bson.M{"$set": userUpdate}

	_, err = usersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return entity.User{}, err
	}

	return s.getUser(id)
}

func (s Store) GetUserForCompany(companyID string) (entity.User, error) {
	usersCollection := s.db.DB.Collection(usersCollectionName)

	filter := filter()

	filter["companyID"] = companyID

	// Sort by createdAt descending
	sort := bson.M{"createdAt": -1}

	var user entity.User
	err := usersCollection.FindOne(context.Background(), filter, &options.FindOneOptions{Sort: sort}).Decode(&user)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s Store) getUser(id string) (entity.User, error) {
	usersCollection := s.db.DB.Collection(usersCollectionName)

	filter := filter()

	filter["_id"] = id

	var user entity.User
	err := usersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func filter() bson.M {
	filter := bson.M{"deletedAt": bson.M{"$exists": false}}

	return filter
}

func randomID() string {
	return uniuri.NewLen(24)
}
