package messageimpl

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/eolso/threadsafe"
	"github.com/worsediscord/server/services/message"
)

type Map struct {
	data *threadsafe.Map[string, *message.Message]
}

func NewMap() *Map {
	return &Map{
		data: threadsafe.NewMap[string, *message.Message](),
	}
}

func (m *Map) Create(_ context.Context, opts message.CreateMessageOpts) error {
	id := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%d", opts.RoomId, time.Now().UnixMilli())))

	msg := message.Message{
		Id:        id,
		UserId:    opts.UserId,
		RoomId:    opts.RoomId,
		Content:   opts.Content,
		Timestamp: time.Now().UnixMilli(),
	}

	m.data.Set(id, &msg)

	return nil
}

func (m *Map) GetMessageById(_ context.Context, opts message.GetMessageByIdOpts) (*message.Message, error) {
	msg, ok := m.data.Get(opts.Id)
	if !ok {
		return nil, message.ErrNotFound
	}

	return msg, nil
}

func (m *Map) List(_ context.Context, _ message.ListMessageOpts) ([]*message.Message, error) {
	return m.data.Values(), nil
}
