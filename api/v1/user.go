package v1

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID           string `json:"username"`
	Password     string `json:"password"`
	Rooms        []Room `json:"rooms"`
	CreationTime string `json:"creation_time"`
	LastActivity string `json:"last_activity"`

	apikey string
}

type UserManager struct {
	ActiveUsers map[string]User
	ErrChan     <-chan error

	lock      sync.Mutex // TODO replace this with something less dumb
	flusher   Flusher
	flushChan chan<- interface{}
}

// NewUser returns a User object using the name provided.
func NewUser(userID string) User {
	return User{
		ID:           userID,
		Rooms:        nil,
		CreationTime: time.Now().Format(time.RFC822Z),
		LastActivity: time.Now().Format(time.RFC822Z),
	}
}

func (u User) WithPassword(password string) User {
	if password != "" {
		u.Password = base64.StdEncoding.EncodeToString([]byte(password))
	}

	return u
}

func NewUserManager() *UserManager {
	return &UserManager{
		ActiveUsers: make(map[string]User),
	}
}

func (um *UserManager) WithFlusher(ctx context.Context, f Flusher) *UserManager {
	um.flusher = f
	um.flushChan, um.ErrChan = f.Listen(ctx)

	return um
}

func (um *UserManager) AddUser(user User) error {
	um.lock.Lock()
	defer um.lock.Unlock()

	if user.ID == "" {
		return fmt.Errorf("user must contain ID: invalid user")
	}

	if _, ok := um.ActiveUsers[user.ID]; ok {
		return fmt.Errorf("user already exists")
	}

	um.ActiveUsers[user.ID] = user

	return nil
}

func (um *UserManager) GetUserByID(userID string) (User, bool) {
	um.lock.Lock()
	defer um.lock.Unlock()

	if _, ok := um.ActiveUsers[userID]; ok {
		return um.ActiveUsers[userID], true
	}

	return User{}, false
}

func (um *UserManager) ListUsers() []User {
	um.lock.Lock()
	defer um.lock.Unlock()

	var users []User
	for k, _ := range um.ActiveUsers {
		users = append(users, um.ActiveUsers[k])
	}

	return users
}

func (um *UserManager) DeleteUser(userID string) error {
	um.lock.Lock()
	defer um.lock.Unlock()

	if userID == "" {
		return nil
	}

	if _, ok := um.ActiveUsers[userID]; !ok {
		return fmt.Errorf("room does not exist")
	}

	delete(um.ActiveUsers, userID)

	return nil
}

func (um *UserManager) CredMap() map[string]string {
	um.lock.Lock()
	defer um.lock.Unlock()

	credMap := make(map[string]string)

	// TODO should probably key on something unique lol
	for _, u := range um.ActiveUsers {
		password, _ := base64.StdEncoding.DecodeString(u.Password)
		credMap[u.ID] = string(password)
	}

	return credMap
}

func (um *UserManager) Save(userID string) error {
	um.lock.Lock()
	defer um.lock.Unlock()

	if userID == "" {
		return nil
	}

	if _, ok := um.ActiveUsers[userID]; !ok {
		return fmt.Errorf("user does not exist")
	}

	um.flushChan <- um.flusher.Flatten(userID, um.ActiveUsers[userID])
	delete(um.ActiveUsers, userID)

	return nil
}

func (um *UserManager) Flush() error {
	um.lock.Lock()
	defer um.lock.Unlock()

	if um.flushChan == nil {
		um.ActiveUsers = nil
		return nil
	}

	for _, u := range um.ActiveUsers {
		um.flushChan <- um.flusher.Flatten(u.ID, u)
	}

	um.ActiveUsers = nil

	return nil
}
