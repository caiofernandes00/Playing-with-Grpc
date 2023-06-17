package repository

import (
	"sync"

	"github.com/caiofernandes00/playing-with-golang/grpc/internal/entity"
)

type UserStore interface {
	Save(user *entity.User) error
	Find(username string) (*entity.User, error)
}

type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*entity.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*entity.User),
	}
}

func (store *InMemoryUserStore) Save(user *entity.User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, ok := store.users[user.Username]; ok {
		return ErrAlreadyExists
	}

	store.users[user.Username] = user.Clone()
	return nil
}

func (store *InMemoryUserStore) Find(username string) (*entity.User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user, ok := store.users[username]
	if !ok {
		return nil, nil
	}

	return user.Clone(), nil
}
