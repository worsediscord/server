package v1

import (
	"encoding/base64"
	"sync"
	"time"
)

type UserMetadata struct {
	ID           string `json:"username"`
	Password     string `json:"-"`
	CreationTime string `json:"creation_time"`
	LastActivity string `json:"last_activity"`
}
type User struct {
	Metadata UserMetadata `json:"metadata"`
	Rooms    []Identity   `json:"-"`

	apikey   string
	roomLock sync.RWMutex
}

// NewUser returns a User object using the name provided.
func NewUser(userID string) *User {
	return &User{
		Metadata: UserMetadata{
			ID:           userID,
			CreationTime: time.Now().Format(time.RFC822Z),
			LastActivity: time.Now().Format(time.RFC822Z),
		},
		Rooms: nil,
	}
}

func BuildUser(md UserMetadata, rooms []Identity) *User {
	return &User{
		Metadata: md,
		Rooms:    rooms,
	}
}

func (u *User) WithPassword(password string) *User {
	if password != "" {
		u.Metadata.Password = base64.StdEncoding.EncodeToString([]byte(password))
		u.updateLastActivity()
	}

	return u
}

func (u *User) UpdatePassword(password string) {
	if password != "" {
		u.Metadata.Password = base64.StdEncoding.EncodeToString([]byte(password))
		u.updateLastActivity()
	}
}

func (u *User) AddRoom(room Identifiable) {
	u.roomLock.Lock()
	defer u.roomLock.Unlock()

	for _, r := range u.Rooms {
		if r.ID == room.UID() {
			return
		}
	}

	u.Rooms = append(u.Rooms, Identity{Name: room.CommonName(), ID: room.UID()})
}

func (u *User) updateLastActivity() {
	u.Metadata.LastActivity = time.Now().Format(time.RFC822Z)
}
