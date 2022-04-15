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

func NewMessage(text string, id Identifiable) Message {
	guid := xid.New()

	return Message{
		ID:        guid.String(),
		Text:      text,
		AuthorID:  id.UID(),
		Author:    id.CommonName(),
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
