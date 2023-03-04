package v1

import (
	"github.com/rs/xid"
	"time"
)

type Room struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	OwnerID      string `json:"owner_id"`
	CreationTime string `json:"creation_time"`
	LastActivity string `json:"last_activity"`

	//Messages       []Message `json:"messages"`
	//PendingInvites []string  `json:"-"`

	//msgLock sync.RWMutex
	//invLock sync.RWMutex
}

type RoomInvite struct {
	Invitee      string `json:"invitee"`
	Invitor      string `json:"invitor"`
	CreationTime string `json:"creation_time"`
	ExpiresAt    string `json:"expires_at"`
}

// NewUser returns a User object using the name provided.
func NewRoom(name string) *Room {
	guid := xid.New()
	return &Room{
		ID:           guid.String(),
		Name:         name,
		CreationTime: time.Now().Format(time.RFC3339),
		LastActivity: time.Now().Format(time.RFC3339),
		//Messages:     nil,
		//PendingInvites: nil,
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

func (r *Room) updateLastActivity() {
	r.LastActivity = time.Now().Format(time.RFC822Z)
}

func NewRoomInvite(invitee string, invitor string, duration time.Duration) RoomInvite {
	return RoomInvite{
		Invitee:      invitee,
		Invitor:      invitor,
		CreationTime: time.Now().Format(time.RFC3339),
		ExpiresAt:    time.Now().Add(duration).Format(time.RFC3339),
	}
}
