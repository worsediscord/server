package v1

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/rs/xid"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	UserCreate = "CREATE"
	UserDelete = "DELETE"
)

type UserRequest struct {
	Method string
	User   User
}

type User struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	LastActive time.Time `json:"last_active"`
}

type UserManager struct {
	activeUsers []User
	userPath    string
	basePath    string
}

func CreateUserHandler(userChan chan<- UserRequest) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest)
			return
		}

		var u User
		err = json.Unmarshal(b, &u)
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest)
			return
		}

		// Validate the request body
		if u.Name == "" || u.Password == "" {
			_ = render.Render(w, r, ErrBadRequest)
			return
		}

		guid := xid.New()
		u.ID = guid.String()
		u.Password = base64.StdEncoding.EncodeToString([]byte(u.Password))

		userChan <- UserRequest{Method: UserCreate, User: u}

		w.WriteHeader(http.StatusOK)
	}
}

func ListUserHandler(um *UserManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//for _, u := um.activeUsers {
		//
		//}
	}
}

func (u User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("required field missing: name")
	} else if u.ID == "" {
		return fmt.Errorf("required field missing: id")
	} else if u.Password == "" {
		return fmt.Errorf("required field missing: password")
	}

	return nil
}

// NewUserManager creates
func NewUserManager(p string) *UserManager {
	return &UserManager{
		basePath: p,
		userPath: filepath.Join(p, "users"),
	}
}

func (um *UserManager) Serve() (chan<- UserRequest, <-chan error) {
	urChan := make(chan UserRequest, 10)
	errChan := make(chan error, 10)

	go func() {
		for ur := range urChan {
			switch ur.Method {
			case UserCreate:
				if err := um.Create(ur.User); err != nil {
					errChan <- err
				}
			case UserDelete:
				if err := um.Delete(ur.User); err != nil {
					errChan <- err
				}
			default:
				// do something
			}
		}
	}()

	return urChan, errChan
}

func (um *UserManager) Create(u User) error {
	for _, user := range um.ListIDs() {
		if user == u.ID {
			return fmt.Errorf("user already exists")
		}
	}

	um.activeUsers = append(um.activeUsers, u)

	return nil
}

func (um *UserManager) Delete(u User) error {
	return nil
}

func (um *UserManager) Load(id string) error {
	path := filepath.Join(um.userPath, id)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("user not found")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open user file")
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read user file: %w", err)
	}

	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %w", err)
	}

	if err = u.Validate(); err != nil {
		return err
	}

	um.activeUsers = append(um.activeUsers, u)

	return nil
}

// Save writes the specified user to disk and removes them from the active user list.
func (um *UserManager) Save(id string) error {
	for _, u := range um.activeUsers {
		if u.ID == id {
			b, err := json.Marshal(u)
			if err != nil {
				return err
			}

			return os.WriteFile(filepath.Join(um.userPath, u.ID), b, 0600)
		}
	}

	return fmt.Errorf("user not found")
}

func (um *UserManager) Flush() error {
	err := os.MkdirAll(um.userPath, 0700)
	if err != nil {
		return err
	}

	for _, u := range um.activeUsers {
		path := filepath.Join(um.userPath, u.ID)

		b, err := json.MarshalIndent(u, "", "  ")
		if err != nil {
			return err
		}

		f, err := os.Create(path)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// ListIDs returns a list of all existing user IDs. This includes active and inactive.
func (um *UserManager) ListIDs() []string {
	var mergedUsers []string
	var userMap map[string]interface{}

	_ = filepath.WalkDir(um.userPath, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if de.IsDir() {
			return nil
		}

		userMap[path] = nil

		return nil
	})

	for _, u := range um.activeUsers {
		userMap[u.ID] = nil
	}

	for k, _ := range userMap {
		mergedUsers = append(mergedUsers, k)
	}

	return mergedUsers
}
