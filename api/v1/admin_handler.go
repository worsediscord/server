package v1

// AdminListRoomHandler lists all rooms ever. Admin only.
//func AdminListRoomHandler(roomsCollection *memcache.Collection) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		response := struct {
//			Rooms []Room `json:"rooms"`
//		}{}
//
//		roomDocs := roomCollection.GetAll()
//		for _, roomDoc := range roomDocs {
//			var room Room
//			if err := roomDoc.Get("metadata").Decode(&room); err != nil {
//				w.WriteHeader(http.StatusInternalServerError)
//				return
//			}
//			response.Rooms = append(response.Rooms, room)
//		}
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
