package v1

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type RoomManager struct {
	activeRooms []RoomState
	workers     int8
	basePath    string
	roomPath    string
}

type CreateRoomBody struct {
	Name string `json:"name"`
}

type RoomMember struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	JoinDate     time.Time `json:"join_date"`
	LastActivity time.Time `json:"last_activity"`
}

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Author    string    `json:"author"`
	Message   string    `json:"message"`
}

type RoomState struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	OwnerID  string       `json:"owner_id"`
	Members  []RoomMember `json:"members"`
	Messages []Message    `json:"messages"`
}

func CreateRoomHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, ok := r.BasicAuth()
		if !ok {
			_ = render.Render(w, r, ErrUnauthorized)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest)
			return
		}

		var crb CreateRoomBody
		err = json.Unmarshal(b, &crb)
		if err != nil || crb.Validate() != nil {
			_ = render.Render(w, r, ErrBadRequest)
			return
		}

		err = rm.Create(crb.Name, user)
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func RemoveRoomHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ID")
		if id == "" {
			_ = render.Render(w, r, ErrBadRequest)
			return
		}

		_, _, ok := r.BasicAuth()
		if !ok {
			_ = render.Render(w, r, ErrUnauthorized)
			return
		}

		for i, room := range rm.activeRooms {
			if room.ID == id {
				rm.activeRooms = append(rm.activeRooms[:i], rm.activeRooms[i+1:]...)
				w.WriteHeader(200)
			}
		}
	}
}

func ListRoomsHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, ok := r.BasicAuth()
		if !ok {
			_ = render.Render(w, r, ErrUnauthorized)
			return
		}

		listRoomsResponse := struct {
			//Rooms []string `json:"rooms"`
			Rooms []struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"rooms"`
		}{}

		for _, room := range rm.activeRooms {
			for _, m := range room.Members {
				if m.Name == user {
					resp := struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					}{
						Name: room.Name,
						ID:   room.ID,
					}
					listRoomsResponse.Rooms = append(listRoomsResponse.Rooms, resp)
				}
			}
		}

		b, err := json.Marshal(listRoomsResponse)
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest)
		}

		w.Write(b)
	}
}

func GetRoomHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ID")
		if id == "" {
			_ = render.Render(w, r, ErrBadRequest)
			return
		}

		_, _, ok := r.BasicAuth()
		if !ok {
			_ = render.Render(w, r, ErrUnauthorized)
			return
		}

		for _, room := range rm.activeRooms {
			if room.ID == id {
				b, err := json.Marshal(room)
				if err != nil {
					_ = render.Render(w, r, ErrBadRequest)
					return
				}
				w.Write(b)
				w.WriteHeader(200)
			}
		}
	}
}

func InviteRoomHandler(rm *RoomManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func NewRoomManager(p string) *RoomManager {
	return &RoomManager{
		basePath: p,
		roomPath: filepath.Join(p, "rooms"),
	}
}

func (rm *RoomManager) Create(name string, owner string) error {
	guid := xid.New()

	rs := RoomState{
		ID:      guid.String(),
		Name:    name,
		OwnerID: owner,
		Members: []RoomMember{
			{
				Name:         owner,
				JoinDate:     time.Now(),
				LastActivity: time.Now(),
			},
		},
	}

	rm.activeRooms = append(rm.activeRooms, rs)

	return nil
}

func (rm *RoomManager) Flush() error {
	err := os.MkdirAll(rm.roomPath, 0700)
	if err != nil {
		return err
	}

	for _, r := range rm.activeRooms {
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(rm.roomPath, r.ID))
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

func (crb CreateRoomBody) Validate() error {
	if crb.Name == "" {
		return fmt.Errorf("required field missing: name")
	}

	return nil
}
