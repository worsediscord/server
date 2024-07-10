package message

type CreateMessageOpts struct {
	UserId  string
	RoomId  int64
	Content string
}

type GetMessageByIdOpts struct {
	Id string
}

type ListMessageOpts struct {
	UserId string
	RoomId int64
}
