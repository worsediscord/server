package v2

import (
	"encoding/base64"
	"time"
)

type User struct {
	ID           string `json:"username"`
	Password     string `json:"-"`
	CreationTime string `json:"creation_time"`
	LastActivity string `json:"last_activity"`

	apikey string
}

// NewUser returns a User object using the name provided.
func NewUser(userID string) *User {
	return &User{
		ID:           userID,
		CreationTime: time.Now().Format(time.RFC822Z),
		LastActivity: time.Now().Format(time.RFC822Z),
	}
}

func (u *User) WithPassword(password string) *User {
	if password != "" {
		u.Password = base64.StdEncoding.EncodeToString([]byte(password))
		u.updateLastActivity()
	}

	return u
}

func (u *User) UpdatePassword(password string) {
	if password != "" {
		u.Password = base64.StdEncoding.EncodeToString([]byte(password))
		u.updateLastActivity()
	}
}

func (u *User) updateLastActivity() {
	u.LastActivity = time.Now().Format(time.RFC822Z)
}
