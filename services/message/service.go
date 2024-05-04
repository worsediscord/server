package message

import "context"

type Service interface {
	Create(context.Context, CreateMessageOpts) error
	GetMessageById(context.Context, GetMessageByIdOpts) (*Message, error)
	List(context.Context, ListMessageOpts) ([]*Message, error)
}
