package v1

import (
	"github.com/rs/xid"
	"time"
)

type Message struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	AuthorID  string `json:"author_id"`
	Author    string `json:"author"`
	Timestamp string `json:"timestamp"`
}

//type MessageList struct {
//	Messages []Message `json:"messages"`
//	lock sync.RWMutex
//}

func NewMessage(text string, author Identifiable) Message {
	guid := xid.New()

	return Message{
		ID:        guid.String(),
		Text:      text,
		AuthorID:  author.UID(),
		Author:    author.CommonName(),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

//
//func (ml *MessageList) SendMessage(message Message) {
//	ml.lock.Lock()
//	ml.Messages = append(ml.Messages, message)
//	ml.lock.Unlock()
//}
//
//func (ml *MessageList) SendSystemMessage(message string) {
//	systemID := Identity{Name: "system", ID: "system"}
//	m := NewMessage(message, systemID)
//	ml.SendMessage(m)
//}
//
//func (ml *MessageList) ListMessages() []Message {
//	ml.lock.RLock()
//	defer ml.lock.RUnlock()
//	return ml.Messages
//}
//
//func (ml *MessageList) DeleteMessage(messageID string) {
//	ml.lock.Lock()
//	defer ml.lock.Unlock()
//
//	for i, m := range ml.Messages {
//		if m.ID == messageID {
//			ml.Messages = append(ml.Messages[:i], ml.Messages[i+1:]...)
//		}
//	}
//}
