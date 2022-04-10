package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/rs/xid"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
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
	urChan := make(chan UserRequest)
	errChan := make(chan error)

	go func() {
		for ur := range urChan {
			switch ur.Method {
			case UserCreate:
				fmt.Println("CREATING")
				if err := um.Create(ur.User); err != nil {
					errChan <- err
				}
			case UserDelete:
				fmt.Println("DELETING")
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
	//b, err := json.Marshal(u)
	//if err != nil {
	//	return err
	//}
	//
	//return os.WriteFile(filepath.Join(um.userPath, u.ID), b, 0600)
}

func (um *UserManager) Delete(u User) error {
	return nil
}

func (um *UserManager) Flush() error {
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
