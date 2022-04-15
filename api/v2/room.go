package v2

import (
	"context"
	"fmt"
	"github.com/rs/xid"
	"sync"
	"time"
)

type Room struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Users        []*User    `json:"users"`
	Messages     []*Message `json:"messages"`
	OwnerID      string     `json:"owner_id"`
	CreationTime string     `json:"creation_time"`
	LastActivity string     `json:"last_activity"`
}

type RoomManager struct {
	ActiveRooms map[string]*Room
	ErrChan     <-chan error

	lock      sync.Mutex // TODO replace this with something less dumb
	flusher   Flusher
	flushChan chan<- interface{}
}

// NewUser returns a User object using the name provided.
func NewRoom(name string) *Room {
	guid := xid.New()
	return &Room{
		ID:           guid.String(),
		Name:         name,
		Users:        nil,
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

func (r *Room) AddUser(user *User) error {
	for _, u := range r.Users {
		if u.ID == user.ID {
			return fmt.Errorf("user %s already exists in room %s", user.ID, r.ID)
		}
	}

	r.Users = append(r.Users, user)

	return nil
}

func (r *Room) GetUser(userID string) (*User, bool) {
	// TODO this should probably be a map, man
	for _, u := range r.Users {
		if u.ID == userID {
			return u, true
		}
	}

	return nil, false
}

func (r *Room) DeleteUser(userID string) error {
	for i, u := range r.Users {
		if u.ID == userID {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("user does not exist")
}

func (r *Room) SendMessage(message *Message) error {
	r.Messages = append(r.Messages, message)

	return nil
}

func (r *Room) ListMessages() []*Message {
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

func NewRoomManager() *RoomManager {
	return &RoomManager{
		ActiveRooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) WithFlusher(ctx context.Context, f Flusher) *RoomManager {
	rm.flusher = f
	rm.flushChan, rm.ErrChan = f.Listen(ctx)

	return rm
}

func (rm *RoomManager) AddRoom(room *Room) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if room.ID == "" {
		return fmt.Errorf("room must contain ID: invalid room")
	}

	if _, ok := rm.ActiveRooms[room.ID]; ok {
		return fmt.Errorf("room already exists")
	}

	rm.ActiveRooms[room.ID] = room

	return nil
}

func (rm *RoomManager) GetRoom(roomID string) (*Room, error) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if roomID == "" {
		return nil, nil // json.Marshal is safe on nil :phew:
	}

	if _, ok := rm.ActiveRooms[roomID]; !ok {
		return nil, fmt.Errorf("room does not exist")
	}

	return rm.ActiveRooms[roomID], nil
}

func (rm *RoomManager) ListRooms() []*Room {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	var rooms []*Room
	for k, _ := range rm.ActiveRooms {
		rooms = append(rooms, rm.ActiveRooms[k])
	}

	return rooms
}

func (rm *RoomManager) DeleteRoom(roomID string) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if roomID == "" {
		return nil
	}

	if _, ok := rm.ActiveRooms[roomID]; !ok {
		return fmt.Errorf("room does not exist")
	}

	delete(rm.ActiveRooms, roomID)

	return nil
}

func (rm *RoomManager) SendMessage(roomID string, message *Message) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if roomID == "" {
		return fmt.Errorf("room must contain ID: invalid room")
	}

	if _, ok := rm.ActiveRooms[roomID]; !ok {
		return fmt.Errorf("room does not exist")
	}

	return rm.ActiveRooms[roomID].SendMessage(message)
}

func (rm *RoomManager) DeleteMessage(roomID string, messageID string) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if roomID == "" {
		return fmt.Errorf("room must contain ID: invalid room")
	}

	if _, ok := rm.ActiveRooms[roomID]; !ok {
		return fmt.Errorf("room does not exist")
	}

	return rm.ActiveRooms[roomID].DeleteMessage(messageID)
}

func (rm *RoomManager) Save(roomID string) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if roomID == "" {
		return nil
	}

	if _, ok := rm.ActiveRooms[roomID]; !ok {
		return fmt.Errorf("room does not exist")
	}

	rm.flushChan <- rm.flusher.Flatten(roomID, rm.ActiveRooms[roomID])
	delete(rm.ActiveRooms, roomID)

	return nil
}

func (rm *RoomManager) Flush() error {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if rm.flushChan == nil {
		rm.ActiveRooms = nil
		return nil
	}

	for _, room := range rm.ActiveRooms {
		rm.flushChan <- rm.flusher.Flatten(room.ID, room)
	}

	rm.ActiveRooms = nil

	return nil
}
