package v2

import (
	"fmt"
	"github.com/rs/xid"
	"sync"
	"time"
)

type Room struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Messages       []Message `json:"messages"`
	OwnerID        string    `json:"owner_id"`
	CreationTime   string    `json:"creation_time"`
	LastActivity   string    `json:"last_activity"`
	PendingInvites []string  `json:"-"`

	msgLock sync.RWMutex
	invLock sync.RWMutex
}

// NewUser returns a User object using the name provided.
func NewRoom(name string) *Room {
	guid := xid.New()
	return &Room{
		ID:           guid.String(),
		Name:         name,
		Messages:     nil,
		CreationTime: time.Now().Format(time.RFC3339),
		LastActivity: time.Now().Format(time.RFC3339),
	}
}

func (r *Room) WithOwner(oid string) *Room {
	r.OwnerID = oid
	return r
}

func (r *Room) CommonName() string {
	return r.Name
}

func (r *Room) UID() string {
	return r.ID
}

func (r *Room) SendMessage(message Message) error {
	r.msgLock.Lock()
	r.Messages = append(r.Messages, message)
	r.msgLock.Unlock()
	r.updateLastActivity()

	return nil
}

func (r *Room) SendSystemMessage(message string) error {
	systemID := Identity{Name: "system", ID: "system"}
	m := NewMessage(message, systemID)

	return r.SendMessage(m)
}

func (r *Room) ListMessages() []Message {
	return r.Messages
}

func (r *Room) DeleteMessage(messageID string) error {
	r.msgLock.RLock()
	for i, m := range r.Messages {
		if m.ID == messageID {
			r.msgLock.RUnlock()
			r.msgLock.Lock()
			r.Messages = append(r.Messages[:i], r.Messages[i+1:]...)
			r.msgLock.Unlock()
			return nil
		}
	}

	r.msgLock.RUnlock()
	return fmt.Errorf("message does not exist")
}

func (r *Room) InviteUser(userID string) {
	r.invLock.RLock()
	for _, inv := range r.PendingInvites {
		if userID == inv {
			r.invLock.RUnlock()
			return
		}
	}
	r.invLock.RUnlock()

	r.invLock.Lock()
	r.PendingInvites = append(r.PendingInvites, userID)
	r.invLock.Unlock()
}

func (r *Room) IsInvited(userID string) bool {
	r.invLock.RLock()
	defer r.invLock.RUnlock()

	for _, inv := range r.PendingInvites {
		if userID == inv {
			return true
		}
	}

	return false
}

func (r *Room) DeleteInvite(userID string) {
	r.invLock.RLock()
	for i, inv := range r.PendingInvites {
		if userID == inv {
			r.invLock.RUnlock()
			r.invLock.Lock()
			r.PendingInvites = append(r.PendingInvites[:i], r.PendingInvites[i+1:]...)
			r.invLock.Unlock()
		}
	}
	r.invLock.RUnlock()
}

func (r *Room) updateLastActivity() {
	r.LastActivity = time.Now().Format(time.RFC822Z)
}
