package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/dchest/uniuri"
	"github.com/noona-hq/blacklist/services/store"
	"github.com/noona-hq/blacklist/services/store/entity"
)

type inMemoryStore struct {
	blacklistUsers map[string][]entity.User // Keyed by companyID
	mu             sync.RWMutex
}

// NewStore creates a new in-memory store.
func NewStore() store.Store {
	return &inMemoryStore{
		blacklistUsers: make(map[string][]entity.User),
	}
}

// CreateBlacklistUser adds a new user to the in-memory store.
func (s *inMemoryStore) CreateBlacklistUser(user entity.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user.ID == "" {
		user.ID = randomID()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Append the user to the slice associated with the company ID.
	s.blacklistUsers[user.CompanyID] = append(s.blacklistUsers[user.CompanyID], user)

	return nil
}

// GetBlacklistUserForCompany retrieves the latest user for a given company ID from the in-memory store.
func (s *inMemoryStore) GetBlacklistUserForCompany(companyID string) (entity.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users, exists := s.blacklistUsers[companyID]
	if !exists || len(users) == 0 {
		return entity.User{}, errors.New("user not found")
	}

	// Since we're always appending, the last user should be the latest one.
	latestUser := users[len(users)-1]

	return latestUser, nil
}

func randomID() string {
	return uniuri.NewLen(24)
}
