package v1

import (
	"encoding/json"
	"github.com/eolso/memcache"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	userMetadataKey = "metadata"
	userRoomsKey    = "rooms"
)

type CreateUserBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserListRoomResponse struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func CreateUserHandler(usersCollection *memcache.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var cub CreateUserBody
		if err = json.Unmarshal(b, &cub); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create the room collection
		thisUserCollection, ok := usersCollection.GetCollection(cub.Name)
		if ok {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			thisUserCollection = usersCollection.Collection(cub.Name)
		}

		user := NewUser(cub.Name).WithPassword(cub.Password)

		// Insert the user metadata into the collection
		thisUserCollection.Document(userMetadataKey).Set("_", user)

		// Insert an empty rooms list into the user's document.
		thisUserCollection.Document(userRoomsKey)

		b, err = json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func GetUserHandler(usersCollection *memcache.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		thisUserCollection, ok := usersCollection.GetCollection(userID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		userDoc, ok := thisUserCollection.GetDocument(userMetadataKey)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		user, _ := userDoc.Get("_")

		b, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func ListUserHandler(usersCollection *memcache.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO do some permission checking against an API key
		//response := struct {
		//	Users []interface{} `json:"users"`
		//}{}

		usersMap := usersCollection.GetCollections()
		for name, documents := range usersMap.Map() {

		}
		//fmt.Printf("%#v\n", usersCollection)
		//fmt.Printf("%#v\n", userDocs.)
		//fmt.Printf("%#v\n", usersCollection.Document("Lemuel17"))
		//
		//fmt.Println(userDocs.Map())
		//for _, doc := range userDocs {
		//	var md UserMetadata
		//	if err := doc.Get("metadata").Decode(&md); err != nil {
		//		w.WriteHeader(http.StatusInternalServerError)
		//		return
		//	}
		//	user := BuildUser(md, nil)
		//	response.Users = append(response.Users, user)
		//}

		b, err := json.Marshal(userDocs.Map())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

//func DeleteUserHandler(userCollection *memcache.Collection) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		userID := chi.URLParam(r, "userID")
//		if userID == "" {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		userDoc := userCollection.Get(userID)
//		if userDoc == nil {
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		userDoc.Delete(userID)
//
//		w.WriteHeader(http.StatusOK)
//	}
//}

//func UserListRoomHandler(userCollection *memcache.Collection, roomCollection *memcache.Collection) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		userID := chi.URLParam(r, "userID")
//		if userID == "" {
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		authedUserID, ok := r.Context().Value("userID").(string)
//		if !ok {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		if userID == "@me" {
//			userID = authedUserID
//		}
//
//		userDoc := userCollection.Get(userID)
//		if userDoc == nil {
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		var rooms []Identity
//		if err := userDoc.Get("rooms").Decode(&rooms); err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		var updated bool
//		for i, room := range rooms {
//			if roomCollection.Get(room.ID) == nil {
//				// The room no longer exists, delete it.
//				rooms = append(rooms[:i], rooms[i+1:]...)
//			}
//		}
//		if updated {
//			if err := userDoc.Set("rooms", rooms); err != nil {
//				w.WriteHeader(http.StatusInternalServerError)
//				return
//			}
//		}
//
//		response := struct {
//			Rooms []Identity `json:"rooms"`
//		}{}
//		response.Rooms = rooms
//
//		b, err := json.Marshal(response)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		w.WriteHeader(http.StatusOK)
//		w.Write(b)
//	}
//}

// LoginUserHandler uses the basic auth header to authenticate a user then returns an API key to use.
func LoginUserHandler(akm *ApiKeyManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, _ := r.BasicAuth()

		apikey, properties := NewApiKey(24, user, time.Hour)

		err := akm.RegisterKey(apikey, properties)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := struct {
			Apikey string `json:"apikey"`
		}{
			Apikey: apikey,
		}

		b, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
