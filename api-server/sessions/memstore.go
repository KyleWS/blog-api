package sessions

import (
	"fmt"
	"time"

	cache "github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
)

type MemStore struct {
	entries *cache.Cache
}

func NewMemStore(sessionDuration time.Duration, purgeInterval time.Duration) *MemStore {
	return &MemStore{
		entries: cache.New(sessionDuration, purgeInterval),
	}
}

func (ms *MemStore) Save(accessToken string, token *oauth2.Token) error {
	ms.entries.Set(accessToken, token, cache.DefaultExpiration)
	return nil
}

func (ms *MemStore) Get(accessToken string) (*oauth2.Token, error) {
	token, present := ms.entries.Get(accessToken)
	if present == false {
		return nil, fmt.Errorf("token not found in database")
	}
	return token.(*oauth2.Token), nil
}

func (ms *MemStore) Delete(accessToken string) error {
	_, err := ms.entries.Get(accessToken)
	if err == false {
		return fmt.Errorf("error deleting token: %v", err)
	}
	ms.entries.Delete(accessToken)
	return nil
}
