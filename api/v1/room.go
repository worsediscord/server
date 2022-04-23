package v2

import (
	"fmt"
	"github.com/rs/xid"
	"time"
)

type RoomUser struct {
	ID          string
	DisplayName string
}

type Room struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Messages       []Message `json:"messages"`
	OwnerID        string    `json:"owner_id"`
	CreationTime   string    `json:"creation_time"`
	LastActivity   string    `json:"last_activity"`
	PendingInvites []string  `json:"-"`
}

func (ru RoomUser) UID() string {
	return ru.ID
}

func (ru RoomUser) CommonName() string {
	return ru.DisplayName
}

// NewUser returns a User object using the name provided.
func NewRoom(name string) Room {
	guid := xid.New()
	return Room{
		ID:           guid.String(),
		Name:         name,
		Messages:     nil,
		CreationTime: time.Now().Format(time.RFC3339),
		LastActivity: time.Now().Format(time.RFC3339),
	}
}

func (r Room) WithOwner(oid string) Room {
	r.OwnerID = oid
	return r
}

func (r Room) CommonName() string {
	return r.Name
}

func (r Room) UID() string {
	return r.ID
}

func (r *Room) SendMessage(message Message) error {
	r.Messages = append(r.Messages, message)

	return nil
}

func (r *Room) ListMessages() []Message {
	return r.Messages
}

func (r *Room) DeleteMessage(messageID string) error {
	for i, m := range r.Messages {
		if m.ID == messageID {
			r.Messages = append(r.Messages[:i], r.Messages[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("message does not exist")
}

func (r *Room) InviteUser(userID string) {
	for _, inv := range r.PendingInvites {
		if userID == inv {
			return
		}
	}

	r.PendingInvites = append(r.PendingInvites, userID)
}

func (r *Room) IsInvited(userID string) bool {
	for _, inv := range r.PendingInvites {
		if userID == inv {
			return true
		}
	}

	return false
}

func (r *Room) DeleteInvite(userID string) {
	for i, inv := range r.PendingInvites {
		if userID == inv {
			r.PendingInvites = append(r.PendingInvites[:i], r.PendingInvites[i+1:]...)
		}
	}
}
