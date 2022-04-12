package v0

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type State struct {
	Messages []Message
	lock     sync.Mutex
}

type Message struct {
	Text      string `json:"text"`
	Author    string `json:"author"`
	Timestamp string `json:"timestamp"`
}

func GetMessageHandler(s *State) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		s.lock.Lock()
		defer s.lock.Unlock()

		f, _ := w.(http.Flusher)

		response := struct {
			Messages []Message `json:"messages"`
		}{}
		response.Messages = s.Messages

		if accept == "application/json" {
			b, _ := json.Marshal(response)
			w.Write(b)
			f.Flush()
		} else {
			for _, m := range s.Messages {
				w.Write([]byte(fmt.Sprintf("[%s] %s: %s\n", m.Timestamp, m.Author, m.Text)))
				f.Flush()
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func SendMessageHandler(s *State) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lock.Lock()
		defer s.lock.Unlock()

		user, _, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		requestBody := struct {
			Message string `json:"message"`
		}{}

		err = json.Unmarshal(b, &requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		m := Message{
			Text:      requestBody.Message,
			Author:    user,
			Timestamp: time.Now().Format(time.RFC1123),
		}

		s.Messages = append(s.Messages, m)

		w.WriteHeader(http.StatusOK)
	}
}
