package message

type Message struct {
	Id        string
	UserId    string
	RoomId    int64
	Content   string
	Timestamp int64
}
