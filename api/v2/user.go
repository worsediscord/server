package v2

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/rs/xid"
	"path/filepath"
	"sync"
	"time"
)

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Rooms        []Room `json:"rooms"`
	CreationTime string `json:"creation_time"`
	LastActivity string `json:"last_activity"`

	apikey string
}

type UserManager struct {
	ActiveUsers map[string]*User
	ErrChan     <-chan error

	lock      sync.Mutex // TODO replace this with something less dumb
	flusher   Flusher
	flushChan chan<- interface{}
}

// NewUser returns a User object using the name provided.
func NewUser(name string) *User {
	guid := xid.New()

	return &User{
		ID:           guid.String(),
		Name:         name,
		Rooms:        nil,
		CreationTime: time.Now().Format(time.RFC822Z),
		LastActivity: time.Now().Format(time.RFC822Z),
	}
}

func (u *User) WithPassword(password string) *User {
	if password != "" {
		u.Password = base64.StdEncoding.EncodeToString([]byte(password))
	}

	return u
}

func (u *User) CommonName() string {
	return u.Name
}

func (u *User) UID() string {
	return u.ID
}

func NewUserManager() *UserManager {
	return &UserManager{
		ActiveUsers: make(map[string]*User),
	}
}

func (um *UserManager) WithFlusher(ctx context.Context, f Flusher) *UserManager {
	um.flusher = f
	um.flushChan, um.ErrChan = f.Listen(ctx)

	return um
}

func (um *UserManager) AddUser(user *User) error {
	um.lock.Lock()
	defer um.lock.Unlock()

	if user.ID == "" {
		return fmt.Errorf("room must contain ID: invalid room")
	}

	if _, ok := um.ActiveUsers[user.ID]; ok {
		return fmt.Errorf("room already exists")
	}

	um.ActiveUsers[user.ID] = user

	return nil
}

// TODO delete this crap because it's a hack work around to current uid vs common name dumbness
func (um *UserManager) LookupUser(cname string) string {
	for k, v := range um.ActiveUsers {
		if v.Name == cname {
			return k
		}
	}

	return ""
}

/////////////////////////////////////////////////////////////////////////////////////////////////

func (um *UserManager) ListUsers() []*User {
	um.lock.Lock()
	defer um.lock.Unlock()

	var users []*User
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
		credMap[u.Name] = string(password)
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
		return fmt.Errorf("room does not exist")
	}

	um.flushChan <- um.flusher.Flatten(userID, um.ActiveUsers[userID])
	delete(um.ActiveUsers, userID)

	return nil
}

func (um *UserManager) Flush() error {
	if um.flushChan == nil {
		um.ActiveUsers = nil
		return nil
	}

	for _, u := range um.ActiveUsers {
		path := filepath.Join("user", u.ID)
		um.flushChan <- um.flusher.Flatten(path, u)
	}

	um.ActiveUsers = nil

	return nil
}
