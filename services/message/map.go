package message

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/eolso/threadsafe"
)

type Map struct {
	data *threadsafe.Map[string, *Message]
}

func NewMap() *Map {
	return &Map{
		data: threadsafe.NewMap[string, *Message](),
	}
}

func (m *Map) Create(_ context.Context, opts CreateMessageOpts) (*Message, error) {
	id := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d%d", opts.RoomId, time.Now().UnixMilli())))

	msg := Message{
		Id:        id,
		UserId:    opts.UserId,
		RoomId:    opts.RoomId,
		Content:   opts.Content,
		Timestamp: time.Now().UnixMilli(),
	}

	m.data.Set(id, &msg)

	return &msg, nil
}

func (m *Map) GetMessageById(_ context.Context, opts GetMessageByIdOpts) (*Message, error) {
	msg, ok := m.data.Get(opts.Id)
	if !ok {
		return nil, ErrNotFound
	}

	return msg, nil
}

func (m *Map) List(_ context.Context, opts ListMessageOpts) ([]*Message, error) {
	messages := make([]*Message, 0)

	for _, msg := range m.data.Values() {
		matchesFilter := true

		if opts.RoomId != 0 && msg.RoomId != opts.RoomId {
			matchesFilter = false
		}
		if len(opts.UserId) > 0 && msg.UserId != opts.UserId {
			matchesFilter = false
		}

		if matchesFilter {
			messages = append(messages, msg)
		}
	}

	return messages, nil
}
