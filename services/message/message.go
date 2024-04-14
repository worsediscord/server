package message

type Message struct {
	Id        string `json:"id,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	RoomId    int64  `json:"room_id,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
